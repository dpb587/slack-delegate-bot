package booland

import (
	"github.com/dpb587/slack-delegate-bot/pkg/condition"
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/message"
)

type Condition struct {
	Conditions []condition.Condition
}

var _ condition.Condition = &Condition{}

func (c Condition) Evaluate(m message.Message) (bool, error) {
	for _, c := range c.Conditions {
		v, err := c.Evaluate(m)
		if err != nil {
			return false, err
		} else if !v {
			return false, nil
		}
	}

	return true, nil
}
