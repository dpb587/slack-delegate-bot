package multiple_test

import (
	"errors"

	"github.com/dpb587/slack-delegate-bot/delegatebot/handler"
	"github.com/dpb587/slack-delegate-bot/delegatebot/handler/handlerfakes"
	. "github.com/dpb587/slack-delegate-bot/delegatebot/handler/handlers/multiple"
	"github.com/dpb587/slack-delegate-bot/delegatebot/message"
	"github.com/dpb587/slack-delegate-bot/logic/delegate"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler", func() {
	var subject Handler
	var firstHandler, secondHandler *handlerfakes.FakeHandler
	var msg message.Message

	BeforeEach(func() {
		firstHandler = &handlerfakes.FakeHandler{}
		firstHandler.ExecuteReturns(
			handler.MessageResponse{
				EmptyMessage: "fake-empty-message",
			},
			nil,
		)

		secondHandler = &handlerfakes.FakeHandler{}
		secondHandler.ExecuteReturns(
			handler.MessageResponse{
				EmptyMessage: "other-fake-empty-message",
			},
			nil,
		)

		subject = Handler{
			Handlers: []handler.Handler{firstHandler}, // , secondHandler},
		}
	})

	Describe("Execute", func() {
		Context("subhandler errors", func() {
			BeforeEach(func() {
				firstHandler.IsApplicableReturns(true, nil)
				firstHandler.ExecuteReturns(handler.MessageResponse{}, errors.New("fake-err1"))
			})

			It("errors", func() {
				_, err := subject.Execute(&msg)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-err1"))
			})
		})

		Context("subhandler has delegates", func() {
			BeforeEach(func() {
				firstHandler.IsApplicableReturns(true, nil)
				firstHandler.ExecuteReturns(
					handler.MessageResponse{
						Delegates: []delegate.Delegate{
							delegate.Literal{Text: "something"},
						},
					},
					nil,
				)
			})

			It("returns the delegates", func() {
				res, err := subject.Execute(&msg)
				Expect(err).NotTo(HaveOccurred())
				Expect(res.Delegates).To(ConsistOf(delegate.Literal{Text: "something"}))
			})
		})

		Context("subhandler has delegates", func() {
			BeforeEach(func() {
				firstHandler.IsApplicableReturns(true, nil)
				firstHandler.ExecuteReturns(
					handler.MessageResponse{
						EmptyMessage: "no delegate available",
					},
					nil,
				)
			})
			It("configures empty message", func() {
				res, err := subject.Execute(&msg)
				Expect(err).NotTo(HaveOccurred())
				Expect(res.Delegates).To(HaveLen(0))
				Expect(res.EmptyMessage).To(Equal("no delegate available"))
			})
		})
	})

	Describe("IsApplicable", func() {
		Context("subhandler applies", func() {
			BeforeEach(func() {
				firstHandler.IsApplicableReturns(true, nil)
			})

			It("applies", func() {
				b, err := subject.IsApplicable(msg)
				Expect(err).NotTo(HaveOccurred())
				Expect(b).To(BeTrue())
			})
		})

		Context("subhandler does not apply", func() {
			BeforeEach(func() {
				firstHandler.IsApplicableReturns(false, nil)
			})

			It("does not apply", func() {
				b, err := subject.IsApplicable(msg)
				Expect(err).NotTo(HaveOccurred())
				Expect(b).To(BeFalse())
			})
		})

		Context("no handlers", func() {
			BeforeEach(func() {
				subject.Handlers = nil
			})

			It("does not apply", func() {
				b, err := subject.IsApplicable(msg)
				Expect(err).NotTo(HaveOccurred())
				Expect(b).To(BeFalse())
			})
		})
	})
})
