package coalesce_test

import (
	"errors"

	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegatefakes"
	. "github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/coalesce"
	"github.com/dpb587/slack-delegate-bot/pkg/message"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Delegator", func() {
	var delegateErr, delegateNone, delegateOne, delegateMany *delegatefakes.FakeDelegator

	BeforeEach(func() {
		delegateErr = &delegatefakes.FakeDelegator{}
		delegateErr.DelegateReturns(nil, errors.New("fake-err1"))

		delegateNone = &delegatefakes.FakeDelegator{}
		delegateNone.DelegateReturns(nil, nil)

		delegateOne = &delegatefakes.FakeDelegator{}
		delegateOne.DelegateReturns([]delegate.Delegate{delegate.Literal{Text: "one"}}, nil)

		delegateMany = &delegatefakes.FakeDelegator{}
		delegateMany.DelegateReturns([]delegate.Delegate{delegate.Literal{Text: "many1"}, delegate.Literal{Text: "many2"}}, nil)
	})

	It("errors early", func() {
		subject := Delegator{
			Delegators: []delegate.Delegator{delegateErr, delegateOne},
		}

		_, err := subject.Delegate(message.Message{})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("fake-err1"))

		Expect(delegateOne.DelegateCallCount()).To(Equal(0))
	})

	It("can return empty", func() {
		subject := Delegator{
			Delegators: []delegate.Delegator{delegateNone},
		}

		found, err := subject.Delegate(message.Message{})
		Expect(err).NotTo(HaveOccurred())
		Expect(found).To(HaveLen(0))
	})

	It("stops with second delegate", func() {
		subject := Delegator{
			Delegators: []delegate.Delegator{delegateNone, delegateOne, delegateMany},
		}

		found, err := subject.Delegate(message.Message{})
		Expect(err).NotTo(HaveOccurred())
		Expect(found).To(ConsistOf(delegate.Literal{Text: "one"}))

		Expect(delegateNone.DelegateCallCount()).To(Equal(1))
		Expect(delegateOne.DelegateCallCount()).To(Equal(1))
		Expect(delegateMany.DelegateCallCount()).To(Equal(0))
	})
})
