package target

import (
	"github.com/dpb587/slack-alias-bot/condition"
	"github.com/dpb587/slack-alias-bot/message"
)

type Condition struct {
	Channel string
}

var _ condition.Condition = &Condition{}

func (c Condition) Evaluate(m message.Message) (bool, error) {
	if m.InterruptTarget == c.Channel {
		return true, nil
	}

	return false, nil
}
