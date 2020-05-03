package slack_test

import (
	"errors"

	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/handler"
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/handler/handlerfakes"
	. "github.com/dpb587/slack-delegate-bot/cmd/delegatebot/service/slack"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
	"github.com/slack-go/slack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageHandler", func() {
	var subject *MessageHandler
	var msg message.Message
	var ev *slack.MessageEvent
	var interruptHandler *handlerfakes.FakeHandler

	BeforeEach(func() {
		interruptHandler = &handlerfakes.FakeHandler{}
		subject = NewMessageHandler(
			slack.New("fake-offline-token").NewRTM(),
			interruptHandler,
		)

		msg = message.Message{
			OriginType: message.ChannelOriginType,
		}

		ev = &slack.MessageEvent{
			Msg: slack.Msg{
				Channel:   "C1234567",
				Timestamp: "fake-timestamp",
			},
		}
	})

	Context("delegate handling", func() {

		It("propagates errors", func() {
			interruptHandler.ExecuteReturns(handler.MessageResponse{}, errors.New("fake-err1"))

			_, err := subject.GetResponse(msg, ev)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err1"))
		})

		Context("delegates provided", func() {
			BeforeEach(func() {
				interruptHandler.ExecuteReturns(
					handler.MessageResponse{
						Delegates: []delegate.Delegate{
							delegate.Literal{Text: "something"},
							delegate.Literal{Text: "completely"},
							delegate.Literal{Text: "different"},
						},
					},
					nil,
				)
			})

			It("responds to direct messages", func() {
				msg.OriginType = message.DirectMessageOriginType

				res, err := subject.GetResponse(msg, ev)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(res.Text).To(Equal("something completely different"))
			})

			It("responds to channels in threads", func() {
				res, err := subject.GetResponse(msg, ev)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(res.Channel).To(Equal("C1234567"))
				Expect(res.ThreadTimestamp).To(Equal("fake-timestamp"))
				Expect(res.Text).To(Equal("^ something completely different"))
			})

			It("responds to existing threads", func() {
				ev.Msg.ThreadTimestamp = "fake-earlier-timestamp"

				res, err := subject.GetResponse(msg, ev)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(res.Channel).To(Equal("C1234567"))
				Expect(res.ThreadTimestamp).To(Equal("fake-earlier-timestamp"))
				Expect(res.Text).To(Equal("^ something completely different"))
			})
		})

		Context("no delegates", func() {
			Context("custom empty messages", func() {
				It("uses it", func() {
					interruptHandler.ExecuteReturns(
						handler.MessageResponse{
							EmptyMessage: "go find your own answer",
						},
						nil,
					)

					res, err := subject.GetResponse(msg, ev)
					Expect(err).NotTo(HaveOccurred())
					Expect(res).ToNot(BeNil())
					Expect(res.Text).To(Equal("go find your own answer"))
				})
			})

			It("stays quiet", func() {
				res, err := subject.GetResponse(msg, ev)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(BeNil())
			})
		})
	})
})
