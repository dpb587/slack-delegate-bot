package topiclookup_test

import (
	"errors"

	"github.com/dpb587/slack-delegate-bot/logic/delegate"
	. "github.com/dpb587/slack-delegate-bot/logic/delegate/delegates/topiclookup"
	"github.com/dpb587/slack-delegate-bot/logic/delegate/delegates/topiclookup/topiclookupfakes"
	"github.com/dpb587/slack-delegate-bot/message"
	"github.com/nlopes/slack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Delegator", func() {
	var fakeSlackAPI *topiclookupfakes.FakeSlackAPI
	var subject Delegator
	var msg message.Message

	BeforeEach(func() {
		fakeSlackAPI = &topiclookupfakes.FakeSlackAPI{}
		subject = Delegator{
			API:     fakeSlackAPI,
			Channel: "C12345678",
		}
	})

	DescribeTable(
		"parsing the real topics",
		func(topic string, expected ...delegate.Delegate) {
			channelInfo := &slack.Channel{}
			channelInfo.Topic = slack.Topic{
				Value: topic,
			}

			fakeSlackAPI.GetChannelInfoReturns(channelInfo, nil)

			actual, err := subject.Delegate(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(ConsistOf(expected))

			Expect(fakeSlackAPI.GetChannelInfoArgsForCall(0)).To(Equal("C12345678"))
		},
		Entry("bbl-users", ":bbl: *interrupt:* <!subteam^S7E4C41HS|@infrastructureteam> note:* bbl _always_ works", delegate.UserGroup{ID: "S7E4C41HS", Alias: "infrastructureteam"}),
		Entry("bbr", "BOSH Backup &amp; Restore | interrupt: <@U08J13EG0> <@UCKK7PZKK> :party_gopher: For PCF/customer specific questions, please ask in the #pcf-backup-restore channel in Pivotal Slack.", delegate.User{ID: "U08J13EG0"}, delegate.User{ID: "UCKK7PZKK"}),
		// Entry("buildpacks", "Interrupt:  `@guillermo` `@ty` `@buildpacks-team` | Lead: `@slevine` | CI: <http://bit.ly/cf-buildpacks|bit.ly/cf-buildpacks> | Java BP: <#C03F5ELTK|java-buildpack> | Hours: 9-6pm EST"),
		Entry("cf-docs", "Questions? Interrupt <@U0JAEKNBH>. Contribute to the Docs! <http://docs.cloudfoundry.org/concepts/contribute.html>", delegate.User{ID: "U0JAEKNBH"}),
		Entry("cli", "Question about Apps or the CC API? Try <#C07C04W4Q|capi> first! Interrupt: <!subteam^S1ZAS8DNY|@cli-team> PM: <@U0CPY3BL2> For contributor discussion, please visit <#CDVP0651P|cli-dev-internal>", delegate.UserGroup{ID: "S1ZAS8DNY", Alias: "cli-team"}),
		Entry("credhub", "Please include your CredHub logs in case of Errors | interrupt: <@U6W2F82B1> <@U8TDZ8VU3> | break glass: `@credhub-team` | PM: <@UDFK4K0KT>, <@UHPMJCXGC>", delegate.User{ID: "U6W2F82B1"}, delegate.User{ID: "U8TDZ8VU3"}),

		// extra, surrounding emojis do not match
		Entry("capi", "Can I push: <http://canibump.cfapps.io|canibump.cfapps.io> Interrupt: :whale: <@U0GQNFF8R> <@U056V1DDK> :boom-avocado:  | PM: <@U91NR3Q3T> :spacewhale2: : | Operators are standing by to take your call 9-6 Pacific"),
	)

	Context("slack errors", func() {
		BeforeEach(func() {
			fakeSlackAPI.GetChannelInfoReturns(nil, errors.New("fake-err1"))
		})

		It("errors", func() {
			_, err := subject.Delegate(msg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err1"))
		})
	})
})
