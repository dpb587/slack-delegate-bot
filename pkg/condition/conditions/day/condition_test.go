package day_test

import (
	"time"

	. "github.com/dpb587/slack-delegate-bot/pkg/condition/conditions/day"
	"github.com/dpb587/slack-delegate-bot/pkg/message"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Condition", func() {
	var subject Condition

	BeforeEach(func() {
		subject = Condition{
			Location: time.UTC,
			Days:     []string{"Mon", "Fri"},
		}
	})

	mustParseRFC3339 := func(value string) time.Time {
		v, err := time.Parse(time.RFC3339, value)
		if err != nil {
			panic(err)
		}

		return v
	}

	Context("non-matching day", func() {
		It("fails", func() {
			b, err := subject.Evaluate(message.Message{Timestamp: mustParseRFC3339("2006-01-03T12:04:05+07:00")})
			Expect(err).NotTo(HaveOccurred())
			Expect(b).To(BeFalse())
		})
	})

	It("succeeds", func() {
		b, err := subject.Evaluate(message.Message{Timestamp: mustParseRFC3339("2006-01-03T03:04:05+07:00")})
		Expect(err).NotTo(HaveOccurred())
		Expect(b).To(BeTrue())
	})
})
