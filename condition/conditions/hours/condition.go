package hours

import (
	"time"

	"github.com/dpb587/slack-delegate-bot/condition"
	"github.com/dpb587/slack-delegate-bot/message"
)

type Condition struct {
	Location *time.Location
	Start    string
	End      string
}

var _ condition.Condition = &Condition{}

func (c Condition) Evaluate(m message.Message) (bool, error) {
	actual := m.Timestamp.In(c.Location).Format("15:04")

	if actual >= c.Start && actual < c.End {
		return true, nil
	}

	return false, nil
}
