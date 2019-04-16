package boolor

import (
	"github.com/dpb587/slack-alias-bot/condition"
	"github.com/dpb587/slack-alias-bot/message"
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
		} else if v {
			return true, nil
		}
	}

	return false, nil
}
