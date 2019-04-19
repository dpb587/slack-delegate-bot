package literalmap_test

import (
	"errors"

	"github.com/dpb587/slack-delegate-bot/interrupt"
	"github.com/dpb587/slack-delegate-bot/interrupt/interruptfakes"
	. "github.com/dpb587/slack-delegate-bot/interrupt/interrupts/literalmap"
	"github.com/dpb587/slack-delegate-bot/message"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interrupt", func() {
	var from *interruptfakes.FakeInterrupt
	var subject Interrupt
	var msg message.Message

	BeforeEach(func() {
		from = &interruptfakes.FakeInterrupt{}
		subject = Interrupt{
			From: from,
			Users: map[string]string{
				"fake-user1": "U12345678",
				"fake-user2": "U23456789",
			},
			Usergroups: map[string]string{
				"fake-usergroup1": "G12345678",
				"fake-usergroup2": "G23456789",
			},
		}
	})

	DescribeTable(
		"parsing the real topics",
		func(input []interrupt.Interruptible, expected []interrupt.Interruptible) {
			from.LookupReturns(input, nil)

			actual, err := subject.Lookup(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(ConsistOf(expected))
		},
		Entry("matched user", []interrupt.Interruptible{interrupt.Literal{Text: "fake-user1"}}, []interrupt.Interruptible{interrupt.User{ID: "U12345678"}}),
		Entry("unknown user", []interrupt.Interruptible{interrupt.Literal{Text: "fake-user3"}}, []interrupt.Interruptible{interrupt.Literal{Text: "fake-user3"}}),
		Entry("matched usergroup", []interrupt.Interruptible{interrupt.Literal{Text: "fake-usergroup1"}}, []interrupt.Interruptible{interrupt.UserGroup{ID: "G12345678"}}),
		Entry("unknown usergroup", []interrupt.Interruptible{interrupt.Literal{Text: "fake-usergroup3"}}, []interrupt.Interruptible{interrupt.Literal{Text: "fake-usergroup3"}}),
		Entry("multiple matches", []interrupt.Interruptible{interrupt.Literal{Text: "fake-user1"}, interrupt.Literal{Text: "fake-usergroup1"}}, []interrupt.Interruptible{interrupt.User{ID: "U12345678"}, interrupt.UserGroup{ID: "G12345678"}}),
		Entry("non-literals are ignored", []interrupt.Interruptible{interrupt.User{ID: "U98765432"}}, []interrupt.Interruptible{interrupt.User{ID: "U98765432"}}),
	)

	It("propagates errors", func() {
		from.LookupReturns(nil, errors.New("fake-err1"))

		_, err := subject.Lookup(msg)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("fake-err1"))
	})
})
