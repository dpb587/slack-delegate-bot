package delegate_test

import (
	. "github.com/dpb587/slack-delegate-bot/pkg/delegate"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserGroup", func() {
	Describe("String", func() {
		It("stringifys", func() {
			Expect(UserGroup{ID: "G12345678", Alias: "fake-name"}.String()).To(Equal("<!subteam^G12345678|@fake-name>"))
		})

		It("stringifys without alias", func() {
			Expect(UserGroup{ID: "G12345678"}.String()).To(Equal("<!subteam^G12345678|@>"))
		})
	})
})
