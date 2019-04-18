package interrupt

import "github.com/dpb587/slack-delegate-bot/message"

type Interrupt interface {
	Lookup(message.Message) ([]Interruptible, error)
}
