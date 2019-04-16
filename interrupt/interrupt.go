package interrupt

import "github.com/dpb587/slack-alias-bot/message"

type Interrupt interface {
	Lookup(message.Message) ([]Interruptible, error)
}
