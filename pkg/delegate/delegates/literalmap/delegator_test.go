package literalmap_test

import (
	"errors"

	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegatefakes"
	. "github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/literalmap"
	"github.com/dpb587/slack-delegate-bot/pkg/message"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Delegator", func() {
	var from *delegatefakes.FakeDelegator
	var subject Delegator
	var msg message.Message

	BeforeEach(func() {
		from = &delegatefakes.FakeDelegator{}
		subject = Delegator{
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
		func(input []delegate.Delegate, expected []delegate.Delegate) {
			from.DelegateReturns(input, nil)

			actual, err := subject.Delegate(msg)
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(ConsistOf(expected))
		},
		Entry("matched user", []delegate.Delegate{delegate.Literal{Text: "fake-user1"}}, []delegate.Delegate{delegate.User{ID: "U12345678"}}),
		Entry("unknown user", []delegate.Delegate{delegate.Literal{Text: "fake-user3"}}, []delegate.Delegate{delegate.Literal{Text: "fake-user3"}}),
		Entry("matched usergroup", []delegate.Delegate{delegate.Literal{Text: "fake-usergroup1"}}, []delegate.Delegate{delegate.UserGroup{ID: "G12345678"}}),
		Entry("unknown usergroup", []delegate.Delegate{delegate.Literal{Text: "fake-usergroup3"}}, []delegate.Delegate{delegate.Literal{Text: "fake-usergroup3"}}),
		Entry("multiple matches", []delegate.Delegate{delegate.Literal{Text: "fake-user1"}, delegate.Literal{Text: "fake-usergroup1"}}, []delegate.Delegate{delegate.User{ID: "U12345678"}, delegate.UserGroup{ID: "G12345678"}}),
		Entry("non-literals are ignored", []delegate.Delegate{delegate.User{ID: "U98765432"}}, []delegate.Delegate{delegate.User{ID: "U98765432"}}),
	)

	It("propagates errors", func() {
		from.DelegateReturns(nil, errors.New("fake-err1"))

		_, err := subject.Delegate(msg)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("fake-err1"))
	})
})
