package usergroup

import (
	"github.com/dpb587/slack-delegate-bot/logic/interrupt"
	"github.com/dpb587/slack-delegate-bot/message"
)

type Interrupt struct {
	ID    string
	Alias string
}

var _ interrupt.Interrupt = &Interrupt{}

func (i Interrupt) Lookup(_ message.Message) ([]interrupt.Interruptible, error) {
	return []interrupt.Interruptible{interrupt.UserGroup{ID: i.ID, Alias: i.Alias}}, nil
}
