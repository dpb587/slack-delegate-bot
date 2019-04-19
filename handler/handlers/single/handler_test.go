package single_test

import (
	"errors"

	"github.com/dpb587/slack-delegate-bot/condition/conditionfakes"
	"github.com/dpb587/slack-delegate-bot/handler"
	. "github.com/dpb587/slack-delegate-bot/handler/handlers/single"
	"github.com/dpb587/slack-delegate-bot/interrupt"
	"github.com/dpb587/slack-delegate-bot/interrupt/interruptfakes"
	"github.com/dpb587/slack-delegate-bot/message"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler", func() {
	var subject Handler
	var msg message.Message
	var int *interruptfakes.FakeInterrupt

	BeforeEach(func() {
		int = &interruptfakes.FakeInterrupt{}
		subject = Handler{
			Interrupt: int,
			Options: handler.Options{
				EmptyMessage: "fake-empty-message",
			},
		}
	})

	Describe("IsApplicable", func() {
		var condition *conditionfakes.FakeCondition

		Context("condition configured", func() {
			BeforeEach(func() {
				condition = &conditionfakes.FakeCondition{}
				subject.Condition = condition
			})

			Context("true", func() {
				BeforeEach(func() {
					condition.EvaluateReturns(true, nil)
				})

				It("applies", func() {
					b, err := subject.IsApplicable(msg)
					Expect(err).NotTo(HaveOccurred())
					Expect(b).To(BeTrue())
				})
			})

			Context("false", func() {
				BeforeEach(func() {
					condition.EvaluateReturns(false, nil)
				})

				It("does not apply", func() {
					b, err := subject.IsApplicable(msg)
					Expect(err).NotTo(HaveOccurred())
					Expect(b).To(BeFalse())
				})
			})
		})

		It("always applies", func() {
			b, err := subject.IsApplicable(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(b).To(BeTrue())
		})
	})

	Describe("Execute", func() {
		Context("interrupt errors", func() {
			BeforeEach(func() {
				int.LookupReturns(nil, errors.New("fake-err1"))
			})

			It("errors", func() {
				_, err := subject.Execute(&msg)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-err1"))
			})
		})

		Context("with interrupts", func() {
			BeforeEach(func() {
				int.LookupReturns([]interrupt.Interruptible{interrupt.Literal{Text: "something"}}, nil)
			})

			It("provides interrupts", func() {
				res, err := subject.Execute(&msg)
				Expect(err).NotTo(HaveOccurred())
				Expect(res.Interrupts).To(ConsistOf(interrupt.Literal{Text: "something"}))
			})
		})

		It("configures empty message", func() {
			res, err := subject.Execute(&msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.Interrupts).To(HaveLen(0))
			Expect(res.EmptyMessage).To(Equal("fake-empty-message"))
		})
	})
})
