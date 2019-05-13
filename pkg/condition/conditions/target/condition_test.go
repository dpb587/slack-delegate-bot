package target_test

import (
	. "github.com/dpb587/slack-delegate-bot/pkg/condition/conditions/target"
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/message"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Condition", func() {
	var subject Condition
	var msg message.Message

	BeforeEach(func() {
		subject = Condition{Channel: "C12345678"}
		msg = message.Message{InterruptTarget: "C12345678"}
	})

	Context("non-matching target", func() {
		BeforeEach(func() {
			msg = message.Message{InterruptTarget: "C98765432"}
		})

		It("fails", func() {
			b, err := subject.Evaluate(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(b).To(BeFalse())
		})
	})

	It("succeeds", func() {
		b, err := subject.Evaluate(msg)
		Expect(err).NotTo(HaveOccurred())
		Expect(b).To(BeTrue())
	})
})
