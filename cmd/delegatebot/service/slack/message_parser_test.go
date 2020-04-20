package slack_test

import (
	"fmt"
	"time"

	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/message"
	. "github.com/dpb587/slack-delegate-bot/cmd/delegatebot/service/slack"
	"github.com/nlopes/slack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageParser", func() {
	var subject *MessageParser
	var msg slack.Msg

	BeforeEach(func() {
		subject = NewMessageParser(&slack.UserDetails{
			ID: "U1234567",
		})
		msg = slack.Msg{
			Channel:   "C1234567",
			Type:      "message",
			Text:      "help me, <@U1234567> you're my only hope.",
			Timestamp: fmt.Sprintf("%d.0", time.Now().Unix()),
		}
	})

	Describe("ParseMessage", func() {
		It("ignores messages from ourselves", func() {
			msg.User = "U1234567"

			res, err := subject.ParseMessage(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(BeNil())
		})

		It("ignores non-message messages", func() {
			msg.Type = "non-message"

			res, err := subject.ParseMessage(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(BeNil())
		})

		It("ignores message deletions", func() {
			msg.SubType = "message_deleted"

			res, err := subject.ParseMessage(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(BeNil())
		})

		It("ignores topic changes", func() {
			msg.SubType = "group_topic"

			res, err := subject.ParseMessage(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(BeNil())
		})

		It("ignores topic change-looking messages", func() {
			msg.Text = "<@U1234567> set the channel topic: something else"

			res, err := subject.ParseMessage(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(BeNil())
		})

		It("parses attachment fallback", func() {
			msg.Attachments = []slack.Attachment{
				{
					Fallback: msg.Text,
				},
			}

			msg.Text = "something else entirely"

			res, err := subject.ParseMessage(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).ToNot(BeNil())
			Expect(res.Text).To(Equal("something else entirely\n\n---\n\nhelp me, <@U1234567> you're my only hope."))
		})

		Context("direct messages", func() {
			BeforeEach(func() {
				msg.Channel = "D1234567"
			})

			Context("channel mentions", func() {
				It("parses it as the interrupt target", func() {
					msg.Text = "<#C3EN0BFC0|credhub>"

					res, err := subject.ParseMessage(msg)
					Expect(err).NotTo(HaveOccurred())
					Expect(res).ToNot(BeNil())
					Expect(res.Origin).To(Equal("D1234567"))
					Expect(res.OriginType).To(Equal(message.DirectMessageOriginType))
					Expect(res.InterruptTarget).To(Equal("C3EN0BFC0"))
				})
			})
		})

		It("ignores messages without a self-mention", func() {
			msg.Text = "chit chat"

			res, err := subject.ParseMessage(msg)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(BeNil())
		})

		It("ignores messages mentioning others", func() {
			msg.Text = "<@U9876543> knows stuff"

			res, err := subject.ParseMessage(msg)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(BeNil())
		})

		It("hears itself mentioned", func() {
			res, err := subject.ParseMessage(msg)
			Expect(err).ToNot(HaveOccurred())
			Expect(res).ToNot(BeNil())
			Expect(res.Origin).To(Equal("C1234567"))
			Expect(res.OriginType).To(Equal(message.ChannelOriginType))
			Expect(res.InterruptTarget).To(Equal("C1234567"))
		})

		Context("external channel reference", func() {
			BeforeEach(func() {
				msg.Text = "hey <#C9876543|star-wars> <@U1234567>, help!"
			})

			It("supports prefixes", func() {
				res, err := subject.ParseMessage(msg)
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(res.Origin).To(Equal("C1234567"))
				Expect(res.OriginType).To(Equal(message.ChannelOriginType))
				Expect(res.InterruptTarget).To(Equal("C9876543"))
			})
		})
	})
})
