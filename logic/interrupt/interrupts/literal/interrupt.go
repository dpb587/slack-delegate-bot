package literal

import (
	"github.com/dpb587/slack-delegate-bot/logic/interrupt"
	"github.com/dpb587/slack-delegate-bot/message"
)

type Interrupt struct {
	Text string
}

var _ interrupt.Interrupt = &Interrupt{}

func (i Interrupt) Lookup(_ message.Message) ([]interrupt.Interruptible, error) {
	return []interrupt.Interruptible{interrupt.Literal{Text: i.Text}}, nil
}
