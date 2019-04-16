package union

import (
	"github.com/dpb587/slack-alias-bot/interrupt"
	"github.com/dpb587/slack-alias-bot/message"
)

type Interrupt struct {
	Interrupts []interrupt.Interrupt
}

var _ interrupt.Interrupt = &Interrupt{}

func (i Interrupt) Lookup(m message.Message) ([]interrupt.Interruptible, error) {
	var r []interrupt.Interruptible

	for _, sub := range i.Interrupts {
		subr, err := sub.Lookup(m)
		if err != nil {
			return nil, err
		}

		r = append(r, subr...)
	}

	return r, nil
}
