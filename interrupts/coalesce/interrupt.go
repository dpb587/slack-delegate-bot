package coalesce

import (
	"github.com/dpb587/slack-alias-bot/interrupt"
	"github.com/dpb587/slack-alias-bot/message"
)

type Interrupt struct {
	Interrupts []interrupt.Interrupt
}

var _ interrupt.Interrupt = &Interrupt{}

func (i Interrupt) Lookup(m message.Message) ([]interrupt.Interruptible, error) {
	for _, sub := range i.Interrupts {
		subr, err := sub.Lookup(m)
		if err != nil {
			return nil, err
		} else if len(subr) > 0 {
			return subr, nil
		}
	}

	return nil, nil
}
