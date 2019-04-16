package boolnot

import (
	"github.com/dpb587/slack-alias-bot/condition"
	"github.com/dpb587/slack-alias-bot/message"
)

type Condition struct {
	Condition condition.Condition
}

var _ condition.Condition = &Condition{}

func (c Condition) Evaluate(m message.Message) (bool, error) {
	v, err := c.Condition.Evaluate(m)

	return !v, err
}
