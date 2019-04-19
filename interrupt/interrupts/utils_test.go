package interrupts_test

import (
	"github.com/dpb587/slack-delegate-bot/interrupt"
	. "github.com/dpb587/slack-delegate-bot/interrupt/interrupts"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
	Describe("Join", func() {
		It("stringifys and joins", func() {
			str := Join(
				[]interrupt.Interruptible{
					interrupt.Literal{Text: "literal"},
					interrupt.User{ID: "U12345678"},
					interrupt.UserGroup{ID: "G12345678"},
				},
				" // ",
			)

			Expect(str).To(Equal("literal // <@U12345678> // <!subteam^G12345678|@>"))
		})
	})
})
