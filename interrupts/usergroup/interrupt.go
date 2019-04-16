package usergroup

import (
	"github.com/dpb587/slack-alias-bot/interrupt"
	"github.com/dpb587/slack-alias-bot/message"
)

type Interrupt struct {
	ID    string
	Alias string
}

var _ interrupt.Interrupt = &Interrupt{}

func (i Interrupt) Lookup(_ message.Message) ([]interrupt.Interruptible, error) {
	return []interrupt.Interruptible{interrupt.UserGroup{ID: i.ID, Alias: i.Alias}}, nil
}
