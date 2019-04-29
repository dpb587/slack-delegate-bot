package delegate_test

import (
	. "github.com/dpb587/slack-delegate-bot/logic/delegate"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User", func() {
	Describe("String", func() {
		It("stringifys", func() {
			Expect(User{ID: "U12345678"}.String()).To(Equal("<@U12345678>"))
		})
	})
})
