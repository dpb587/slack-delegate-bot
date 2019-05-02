package boolor_test

import (
	"errors"

	"github.com/dpb587/slack-delegate-bot/delegatebot/message"
	"github.com/dpb587/slack-delegate-bot/logic/condition"
	"github.com/dpb587/slack-delegate-bot/logic/condition/conditionfakes"
	. "github.com/dpb587/slack-delegate-bot/logic/condition/conditions/boolor"

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
			subject.Conditions = []condition.Condition{falseCondition, trueCondition}

			res, err := subject.Evaluate(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(BeTrue())

			Expect(falseCondition.EvaluateCallCount()).To(Equal(1))
			Expect(trueCondition.EvaluateCallCount()).To(Equal(1))
		})

		It("stops on error", func() {
			subject.Conditions = []condition.Condition{falseCondition, errorCondition, trueCondition}

			_, err := subject.Evaluate(msg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-evaluate-err"))

			Expect(falseCondition.EvaluateCallCount()).To(Equal(1))
			Expect(errorCondition.EvaluateCallCount()).To(Equal(1))
			Expect(trueCondition.EvaluateCallCount()).To(Equal(0))
		})

		It("evalutes to false if nothing matches", func() {
			subject.Conditions = []condition.Condition{falseCondition}

			res, err := subject.Evaluate(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(BeFalse())

			Expect(falseCondition.EvaluateCallCount()).To(Equal(1))
		})
	})
})
