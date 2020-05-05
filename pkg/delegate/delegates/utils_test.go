package delegates_test

import (
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	. "github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/pkg/message"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
	Describe("Join", func() {
		It("stringifys and joins", func() {
			str := Join(
				[]message.Delegate{
					delegate.Literal{Text: "literal"},
					delegate.User{ID: "U12345678"},
					delegate.UserGroup{ID: "G12345678"},
				},
				" // ",
			)

			Expect(str).To(Equal("literal // <@U12345678> // <!subteam^G12345678|@>"))
		})
	})
})
