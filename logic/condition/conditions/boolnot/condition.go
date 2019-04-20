package boolnot

import (
	"github.com/dpb587/slack-delegate-bot/logic/condition"
	"github.com/dpb587/slack-delegate-bot/delegatebot/message"
)

type Condition struct {
	Condition condition.Condition
}

var _ condition.Condition = &Condition{}

func (c Condition) Evaluate(m message.Message) (bool, error) {
	v, err := c.Condition.Evaluate(m)

	return !v, err
}
