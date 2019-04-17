package literal

import (
	"github.com/dpb587/slack-alias-bot/interrupt"
	"github.com/dpb587/slack-alias-bot/message"
)

type Interrupt struct {
	Text string
}

var _ interrupt.Interrupt = &Interrupt{}

func (i Interrupt) Lookup(_ message.Message) ([]interrupt.Interruptible, error) {
	return []interrupt.Interruptible{interrupt.Literal{Text: i.Text}}, nil
}
