package coalesce_test

import (
	"errors"

	"github.com/dpb587/slack-delegate-bot/interrupt"
	"github.com/dpb587/slack-delegate-bot/interrupt/interruptfakes"
	. "github.com/dpb587/slack-delegate-bot/interrupt/interrupts/coalesce"
	"github.com/dpb587/slack-delegate-bot/message"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Interrupt", func() {
	var interruptErr, interruptNone, interruptOne, interruptMany *interruptfakes.FakeInterrupt

	BeforeEach(func() {
		interruptErr = &interruptfakes.FakeInterrupt{}
		interruptErr.LookupReturns(nil, errors.New("fake-err1"))

		interruptNone = &interruptfakes.FakeInterrupt{}
		interruptNone.LookupReturns(nil, nil)

		interruptOne = &interruptfakes.FakeInterrupt{}
		interruptOne.LookupReturns([]interrupt.Interruptible{interrupt.Literal{Text: "one"}}, nil)

		interruptMany = &interruptfakes.FakeInterrupt{}
		interruptMany.LookupReturns([]interrupt.Interruptible{interrupt.Literal{Text: "many1"}, interrupt.Literal{Text: "many2"}}, nil)
	})

	It("errors early", func() {
		subject := Interrupt{
			Interrupts: []interrupt.Interrupt{interruptErr, interruptOne},
		}

		_, err := subject.Lookup(message.Message{})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("fake-err1"))

		Expect(interruptOne.LookupCallCount()).To(Equal(0))
	})

	It("can return empty", func() {
		subject := Interrupt{
			Interrupts: []interrupt.Interrupt{interruptNone},
		}

		found, err := subject.Lookup(message.Message{})
		Expect(err).NotTo(HaveOccurred())
		Expect(found).To(HaveLen(0))
	})

	It("stops with second interrupt", func() {
		subject := Interrupt{
			Interrupts: []interrupt.Interrupt{interruptNone, interruptOne, interruptMany},
		}

		found, err := subject.Lookup(message.Message{})
		Expect(err).NotTo(HaveOccurred())
		Expect(found).To(ConsistOf(interrupt.Literal{Text: "one"}))

		Expect(interruptNone.LookupCallCount()).To(Equal(1))
		Expect(interruptOne.LookupCallCount()).To(Equal(1))
		Expect(interruptMany.LookupCallCount()).To(Equal(0))
	})
})
