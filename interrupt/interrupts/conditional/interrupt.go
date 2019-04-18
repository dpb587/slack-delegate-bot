package conditional

import (
	"github.com/dpb587/slack-delegate-bot/condition"
	"github.com/dpb587/slack-delegate-bot/interrupt"
	"github.com/dpb587/slack-delegate-bot/message"
)

type Interrupt struct {
	When condition.Condition
	Then interrupt.Interrupt
	Else interrupt.Interrupt
}

var _ interrupt.Interrupt = &Interrupt{}

func (i Interrupt) Lookup(m message.Message) ([]interrupt.Interruptible, error) {
	when, err := i.When.Evaluate(m)
	if err != nil {
		return nil, err
	}

	if when {
		if i.Then != nil {
			return i.Then.Lookup(m)
		}
	} else if i.Else != nil {
		return i.Else.Lookup(m)
	}

	return nil, nil
}
