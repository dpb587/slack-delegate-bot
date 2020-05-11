package slack_test

import (
	"errors"

	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegatefakes"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
	. "github.com/dpb587/slack-delegate-bot/pkg/slack"
	"github.com/dpb587/slack-delegate-bot/pkg/slack/slackfakes"
	"github.com/slack-go/slack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageHandler", func() {
	var subject *Responder
	var msg message.Message
	var delegator *delegatefakes.FakeDelegator
	var slackAPI *slackfakes.FakeResponderSlackAPI

	BeforeEach(func() {
		slackAPI = &slackfakes.FakeResponderSlackAPI{}
		delegator = &delegatefakes.FakeDelegator{}

		subject = NewResponder(slackAPI, delegator)

		msg = message.Message{
			ChannelID:    "C1234567",
			RawTimestamp: "fake-timestamp",
			Type:         message.ChannelMessageType,
		}
	})

	Context("delegate handling", func() {
		It("propagates errors", func() {
			delegator.DelegateReturns(nil, errors.New("fake-err1"))

			err := subject.ProcessMessage(msg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err1"))
			Expect(slackAPI.PostMessageCallCount()).To(Equal(0))
		})

		Context("delegates provided", func() {
			BeforeEach(func() {
				delegator.DelegateReturns(
					[]message.Delegate{
						delegate.Literal{Text: "something"},
						delegate.Literal{Text: "completely"},
						delegate.Literal{Text: "different"},
					},
					nil,
				)
			})

			It("responds to direct messages", func() {
				msg.ChannelID = "D1234567"
				msg.Type = message.DirectMessageMessageType

				err := subject.ProcessMessage(msg)
				Expect(err).NotTo(HaveOccurred())
				Expect(slackAPI.PostMessageCallCount()).To(Equal(1))

				channel, opts := slackAPI.PostMessageArgsForCall(0)
				endpoint, values, err := slack.UnsafeApplyMsgOptions("fake-token", channel, "fake-url/", opts...)
				Expect(err).ToNot(HaveOccurred())
				Expect(endpoint).To(Equal("fake-url/chat.postMessage"))
				Expect(values.Get("channel")).To(Equal("D1234567"))
				Expect(values.Get("text")).To(Equal("something completely different"))

				// Expect(res.Text).To(Equal("something completely different"))
			})

			It("responds to channels in threads", func() {
				err := subject.ProcessMessage(msg)
				Expect(err).NotTo(HaveOccurred())
				Expect(slackAPI.PostMessageCallCount()).To(Equal(1))

				channel, opts := slackAPI.PostMessageArgsForCall(0)
				endpoint, values, err := slack.UnsafeApplyMsgOptions("fake-token", channel, "fake-url/", opts...)
				Expect(err).ToNot(HaveOccurred())
				Expect(endpoint).To(Equal("fake-url/chat.postMessage"))
				Expect(values.Get("channel")).To(Equal("C1234567"))
				Expect(values.Get("thread_ts")).To(Equal("fake-timestamp"))
				Expect(values.Get("text")).To(Equal("^ something completely different"))

				// Expect(res).ToNot(BeNil())
				// Expect(res.Channel).To(Equal("C1234567"))
				// Expect(res.ThreadTimestamp).To(Equal("fake-timestamp"))
				// Expect(res.Text).To(Equal("^ something completely different"))
			})

			It("responds to existing threads", func() {
				msg.RawThreadTimestamp = "fake-earlier-timestamp"

				err := subject.ProcessMessage(msg)
				Expect(err).NotTo(HaveOccurred())
				Expect(slackAPI.PostMessageCallCount()).To(Equal(1))

				channel, opts := slackAPI.PostMessageArgsForCall(0)
				endpoint, values, err := slack.UnsafeApplyMsgOptions("fake-token", channel, "fake-url/", opts...)
				Expect(err).ToNot(HaveOccurred())
				Expect(endpoint).To(Equal("fake-url/chat.postMessage"))
				Expect(values.Get("channel")).To(Equal("C1234567"))
				Expect(values.Get("thread_ts")).To(Equal("fake-earlier-timestamp"))
				Expect(values.Get("text")).To(Equal("^ something completely different"))

				// Expect(res).ToNot(BeNil())
				// Expect(res.Channel).To(Equal("C1234567"))
				// Expect(res.ThreadTimestamp).To(Equal("fake-earlier-timestamp"))
				// Expect(res.Text).To(Equal("^ something completely different"))
			})
		})

		Context("no delegates", func() {
			It("stays quiet", func() {
				err := subject.ProcessMessage(msg)
				Expect(err).NotTo(HaveOccurred())
				Expect(slackAPI.PostMessageCallCount()).To(Equal(0))
			})
		})
	})
})
