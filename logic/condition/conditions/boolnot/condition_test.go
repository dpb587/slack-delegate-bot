package boolnot_test

import (
	"github.com/dpb587/slack-delegate-bot/logic/condition/conditionfakes"
	. "github.com/dpb587/slack-delegate-bot/logic/condition/conditions/boolnot"
	"github.com/dpb587/slack-delegate-bot/message"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Condition", func() {
	var condition *conditionfakes.FakeCondition
	var msg message.Message

	BeforeEach(func() {
		condition = &conditionfakes.FakeCondition{}
		condition.EvaluateReturns(true, nil)
	})

	Context("false", func() {
		BeforeEach(func() {
			condition.EvaluateReturns(false, nil)
		})

		It("is true", func() {
			b, err := Condition{Condition: condition}.Evaluate(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(b).To(BeTrue())
		})
	})

	It("inverts true", func() {
		b, err := Condition{Condition: condition}.Evaluate(msg)
		Expect(err).NotTo(HaveOccurred())
		Expect(b).To(BeFalse())
	})
})
