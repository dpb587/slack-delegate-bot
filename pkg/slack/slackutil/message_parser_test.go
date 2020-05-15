package slackutil_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/dpb587/slack-delegate-bot/pkg/message"
	. "github.com/dpb587/slack-delegate-bot/pkg/slack/slackutil"
)

var _ = Describe("MentionParser", func() {
	Describe("ParseMessageForAnyChannelReference", func() {
		It("parses channels", func() {
			msg := ParseMessageForAnyChannelReference(message.Message{
				TargetChannelID: "C1234567",
				RawText:         "tell me about <#C9876543|star-wars>",
			})

			Expect(msg.TargetChannelID).To(Equal("C9876543"))
		})

		It("parses unnamed channels", func() {
			msg := ParseMessageForAnyChannelReference(
				message.Message{
					RawText: "tell me about <#C9876543>",
				},
			)

			Expect(msg.TargetChannelID).To(Equal("C9876543"))
		})
	})

	Describe("ParseMessageForChannelReference", func() {
		DescribeTable(
			"in-context channels",
			func(appUserID, expectedTarget string) {
				msg := ParseMessageForChannelReference(
					message.Message{
						TargetChannelID: "C1234567",
						RawText:         "hey <#C9876543|star-wars> <@U1234567>, help!",
					},
					func(in string) bool {
						return appUserID == in
					},
				)

				Expect(msg.TargetChannelID).To(Equal(expectedTarget))
			},
			Entry("next to app user", "U1234567", "C9876543"),
			Entry("next to random user", "U9876543", "C1234567"),
		)

		It("parses unnamed channels", func() {
			msg := ParseMessageForChannelReference(
				message.Message{
					RawText: "hey <#C9876543> <@U1234567>, help!",
				},
				func(in string) bool {
					return true
				},
			)

			Expect(msg.TargetChannelID).To(Equal("C9876543"))
		})
	})

	Describe("CheckMessageForMention", func() {
		DescribeTable(
			"mentioned user",
			func(appUserID string, expected bool) {
				actual := CheckMessageForMention(
					message.Message{
						TargetChannelID: "C1234567",
						RawText:         "hey <@U1234567>, help!",
					},
					func(in string) bool {
						return appUserID == in
					},
				)

				Expect(actual).To(Equal(expected))
			},
			Entry("with app user", "U1234567", true),
			Entry("with random user", "U9876543", false),
		)
	})
})
