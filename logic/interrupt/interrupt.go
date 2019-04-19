package interrupt

import "github.com/dpb587/slack-delegate-bot/message"

//go:generate counterfeiter . Interrupt
type Interrupt interface {
	Lookup(message.Message) ([]Interruptible, error)
}
