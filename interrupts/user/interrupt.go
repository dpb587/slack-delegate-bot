package user

import (
	"github.com/dpb587/slack-alias-bot/interrupt"
	"github.com/dpb587/slack-alias-bot/message"
)

type Interrupt struct {
	ID string
}

var _ interrupt.Interrupt = &Interrupt{}

func (i Interrupt) Lookup(_ message.Message) ([]interrupt.Interruptible, error) {
	return []interrupt.Interruptible{interrupt.User{ID: i.ID}}, nil
}
