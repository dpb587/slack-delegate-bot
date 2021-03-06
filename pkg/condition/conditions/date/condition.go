package date

import (
	"time"

	"github.com/dpb587/slack-delegate-bot/pkg/condition"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
)

type Condition struct {
	Location *time.Location
	Dates    []string
}

var _ condition.Condition = &Condition{}

func (c Condition) Evaluate(m message.Message) (bool, error) {
	actual := m.Time.In(c.Location).Format("2006-01-02")

	for _, expected := range c.Dates {
		if actual == expected {
			return true, nil
		}
	}

	return false, nil
}
