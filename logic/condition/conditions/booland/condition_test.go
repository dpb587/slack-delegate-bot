package booland_test

import (
	"errors"

	"github.com/dpb587/slack-delegate-bot/delegatebot/message"
	"github.com/dpb587/slack-delegate-bot/logic/condition"
	"github.com/dpb587/slack-delegate-bot/logic/condition/conditionfakes"
	. "github.com/dpb587/slack-delegate-bot/logic/condition/conditions/booland"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Condition", func() {
	var subject Condition
	var msg message.Message

	Describe("Evaluate", func() {
		var errorCondition, falseCondition, trueCondition *conditionfakes.FakeCondition

		BeforeEach(func() {
			subject = Condition{}
			msg = message.Message{}

			errorCondition = &conditionfakes.FakeCondition{}
			errorCondition.EvaluateReturns(false, errors.New("fake-evaluate-err"))

			falseCondition = &conditionfakes.FakeCondition{}

			trueCondition = &conditionfakes.FakeCondition{}
			trueCondition.EvaluateReturns(true, nil)
		})

		It("evaluates multiple conditions", func() {
			subject.Conditions = []condition.Condition{trueCondition, trueCondition}

			res, err := subject.Evaluate(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(BeTrue())

			Expect(trueCondition.EvaluateCallCount()).To(Equal(2))
		})

		It("stops on error", func() {
			subject.Conditions = []condition.Condition{trueCondition, errorCondition, trueCondition}

			_, err := subject.Evaluate(msg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-evaluate-err"))

			Expect(errorCondition.EvaluateCallCount()).To(Equal(1))
			Expect(trueCondition.EvaluateCallCount()).To(Equal(1))
		})

		It("stops on false", func() {
			subject.Conditions = []condition.Condition{trueCondition, falseCondition, trueCondition}

			res, err := subject.Evaluate(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(BeFalse())

			Expect(falseCondition.EvaluateCallCount()).To(Equal(1))
			Expect(trueCondition.EvaluateCallCount()).To(Equal(1))
		})

		It("is true if nothing is false", func() {
			subject.Conditions = []condition.Condition{trueCondition}

			res, err := subject.Evaluate(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(BeTrue())

			Expect(trueCondition.EvaluateCallCount()).To(Equal(1))
		})
	})
})
