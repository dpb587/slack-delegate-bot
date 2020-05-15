package rtm_test

import (
	"fmt"
	"time"

	"github.com/dpb587/slack-delegate-bot/pkg/message"
	. "github.com/dpb587/slack-delegate-bot/pkg/slack/rtm"
	"github.com/slack-go/slack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	var subject *Parser
	var msg slack.Msg

	BeforeEach(func() {
		subject = NewParser(&slack.UserDetails{
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

			_, reply, err := subject.ParseMessage(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(reply).To(BeFalse())
		})

		It("ignores non-message messages", func() {
			msg.Type = "non-message"

			_, reply, err := subject.ParseMessage(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(reply).To(BeFalse())
		})

		It("ignores message deletions", func() {
			msg.SubType = "message_deleted"

			_, reply, err := subject.ParseMessage(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(reply).To(BeFalse())
		})

		It("ignores topic changes", func() {
			msg.SubType = "group_topic"

			_, reply, err := subject.ParseMessage(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(reply).To(BeFalse())
		})

		It("ignores topic change-looking messages", func() {
			msg.Text = "<@U1234567> set the channel topic: something else"

			_, reply, err := subject.ParseMessage(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(reply).To(BeFalse())
		})

		It("parses attachment fallback", func() {
			msg.Attachments = []slack.Attachment{
				{
					Fallback: msg.Text,
				},
			}

			msg.Text = "something else entirely"

			res, reply, err := subject.ParseMessage(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(reply).To(BeTrue())
			Expect(res.RawText).To(Equal("something else entirely\n\n---\n\nhelp me, <@U1234567> you're my only hope."))
		})

		Context("direct messages", func() {
			BeforeEach(func() {
				msg.Channel = "D1234567"
			})

			Context("channel mentions", func() {
				It("parses it as the interrupt target", func() {
					msg.Text = "<#C3EN0BFC0|credhub>"

					res, reply, err := subject.ParseMessage(msg)
					Expect(err).NotTo(HaveOccurred())
					Expect(reply).To(BeTrue())
					Expect(res.ChannelID).To(Equal("D1234567"))
					Expect(res.TargetChannelID).To(Equal("C3EN0BFC0"))
					Expect(res.Type).To(Equal(message.DirectMessageMessageType))
				})
			})
		})

		It("ignores messages without a self-mention", func() {
			msg.Text = "chit chat"

			_, reply, err := subject.ParseMessage(msg)
			Expect(err).ToNot(HaveOccurred())
			Expect(reply).To(BeFalse())
		})

		It("ignores messages mentioning others", func() {
			msg.Text = "<@U9876543> knows stuff"

			_, reply, err := subject.ParseMessage(msg)
			Expect(err).ToNot(HaveOccurred())
			Expect(reply).To(BeFalse())
		})

		It("hears itself mentioned", func() {
			res, reply, err := subject.ParseMessage(msg)
			Expect(err).ToNot(HaveOccurred())
			Expect(reply).To(BeTrue())
			Expect(res.ChannelID).To(Equal("C1234567"))
			Expect(res.TargetChannelID).To(Equal("C1234567"))
			Expect(res.Type).To(Equal(message.ChannelMessageType))
		})

		Context("external channel reference", func() {
			BeforeEach(func() {
				msg.Text = "hey <#C9876543|star-wars> <@U1234567>, help!"
			})

			It("supports prefixes", func() {
				res, reply, err := subject.ParseMessage(msg)
				Expect(err).ToNot(HaveOccurred())
				Expect(reply).To(BeTrue())
				Expect(res.ChannelID).To(Equal("C1234567"))
				Expect(res.TargetChannelID).To(Equal("C9876543"))
				Expect(res.Type).To(Equal(message.ChannelMessageType))
			})
		})
	})
})
