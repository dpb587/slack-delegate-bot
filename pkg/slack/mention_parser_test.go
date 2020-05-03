package slack_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/dpb587/slack-delegate-bot/pkg/message"
	. "github.com/dpb587/slack-delegate-bot/pkg/slack"
)

var _ = Describe("MentionParser", func() {
	Describe("ParseMessageForAnyChannelReference", func() {
		It("parses channels", func() {
			msg := ParseMessageForAnyChannelReference(message.Message{
				InterruptTarget: "C1234567",
				Text:            "tell me about <#C9876543|star-wars>",
			})

			Expect(msg.InterruptTarget).To(Equal("C9876543"))
		})

		It("parses unnamed channels", func() {
			msg := ParseMessageForAnyChannelReference(
				message.Message{
					Text: "tell me about <#C9876543>",
				},
			)

			Expect(msg.InterruptTarget).To(Equal("C9876543"))
		})
	})

	Describe("ParseMessageForChannelReference", func() {
		DescribeTable(
			"in-context channels",
			func(appUserID, expectedTarget string) {
				msg := ParseMessageForChannelReference(
					message.Message{
						InterruptTarget: "C1234567",
						Text:            "hey <#C9876543|star-wars> <@U1234567>, help!",
					},
					func(in string) bool {
						return appUserID == in
					},
				)

				Expect(msg.InterruptTarget).To(Equal(expectedTarget))
			},
			Entry("next to app user", "U1234567", "C9876543"),
			Entry("next to random user", "U9876543", "C1234567"),
		)

		It("parses unnamed channels", func() {
			msg := ParseMessageForChannelReference(
				message.Message{
					Text: "hey <#C9876543> <@U1234567>, help!",
				},
				func(in string) bool {
					return true
				},
			)

			Expect(msg.InterruptTarget).To(Equal("C9876543"))
		})
	})

	Describe("CheckMessageForMention", func() {
		DescribeTable(
			"mentioned user",
			func(appUserID string, expected bool) {
				actual := CheckMessageForMention(
					message.Message{
						InterruptTarget: "C1234567",
						Text:            "hey <@U1234567>, help!",
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
