package topiclookup_test

import (
	"github.com/dpb587/slack-delegate-bot/interrupt"
	. "github.com/dpb587/slack-delegate-bot/interrupt/interrupts/topiclookup"
	"github.com/dpb587/slack-delegate-bot/interrupt/interrupts/topiclookup/topiclookupfakes"
	"github.com/dpb587/slack-delegate-bot/message"
	"github.com/nlopes/slack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interrupt", func() {
	DescribeTable(
		"the parsing",
		func(topic string, expected ...interrupt.Interruptible) {
			channelInfo := &slack.Channel{}
			channelInfo.Topic = slack.Topic{
				Value: topic,
			}

			fakeSlackAPI := &topiclookupfakes.FakeSlackAPI{}
			fakeSlackAPI.GetChannelInfoReturns(channelInfo, nil)

			subject := Interrupt{
				API:     fakeSlackAPI,
				Channel: "C12345678",
			}

			actual, err := subject.Lookup(message.Message{})
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(ConsistOf(expected))

			Expect(fakeSlackAPI.GetChannelInfoArgsForCall(0)).To(Equal("C12345678"))
		},
		Entry("bbl-users", ":bbl: *interrupt:* <!subteam^S7E4C41HS|@infrastructureteam> note:* bbl _always_ works", interrupt.UserGroup{ID: "S7E4C41HS", Alias: "infrastructureteam"}),
		Entry("bbr", "BOSH Backup &amp; Restore | interrupt: <@U08J13EG0> <@UCKK7PZKK> :party_gopher: For PCF/customer specific questions, please ask in the #pcf-backup-restore channel in Pivotal Slack.", interrupt.User{ID: "U08J13EG0"}, interrupt.User{ID: "UCKK7PZKK"}),
		// Entry("buildpacks", "Interrupt:  `@guillermo` `@ty` `@buildpacks-team` | Lead: `@slevine` | CI: <http://bit.ly/cf-buildpacks|bit.ly/cf-buildpacks> | Java BP: <#C03F5ELTK|java-buildpack> | Hours: 9-6pm EST"),
		// Entry("capi", "Can I push: <http://canibump.cfapps.io|canibump.cfapps.io> Interrupt: :whale: <@U0GQNFF8R> <@U056V1DDK> :boom-avocado:  | PM: <@U91NR3Q3T> :spacewhale2: : | Operators are standing by to take your call 9-6 Pacific", interrupt.User{ID: "U0GQNFF8R"}, interrupt.User{ID: "U056V1DDK"}),
		Entry("cf-docs", "Questions? Interrupt <@U0JAEKNBH>. Contribute to the Docs! <http://docs.cloudfoundry.org/concepts/contribute.html>", interrupt.User{ID: "U0JAEKNBH"}),
		Entry("cli", "Question about Apps or the CC API? Try <#C07C04W4Q|capi> first! Interrupt: <!subteam^S1ZAS8DNY|@cli-team> PM: <@U0CPY3BL2> For contributor discussion, please visit <#CDVP0651P|cli-dev-internal>", interrupt.UserGroup{ID: "S1ZAS8DNY", Alias: "cli-team"}),
		Entry("credhub", "Please include your CredHub logs in case of Errors | interrupt: <@U6W2F82B1> <@U8TDZ8VU3> | break glass: `@credhub-team` | PM: <@UDFK4K0KT>, <@UHPMJCXGC>", interrupt.User{ID: "U6W2F82B1"}, interrupt.User{ID: "U8TDZ8VU3"}),
	)
})
