package delegate_test

import (
	. "github.com/dpb587/slack-delegate-bot/pkg/delegate"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Literal", func() {
	Describe("String", func() {
		It("is very literal", func() {
			Expect(Literal{Text: "one two"}.String()).To(Equal("one two"))
		})
	})
})
