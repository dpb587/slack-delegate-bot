package target

import (
	"github.com/dpb587/slack-delegate-bot/logic/condition"
	"github.com/dpb587/slack-delegate-bot/delegatebot/message"
)

type Condition struct {
	Channel string
}

var _ condition.Condition = &Condition{}

func (c Condition) Evaluate(m message.Message) (bool, error) {
	return m.InterruptTarget == c.Channel, nil
}
