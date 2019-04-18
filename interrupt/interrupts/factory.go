package interrupts

import "github.com/dpb587/slack-delegate-bot/interrupt"

type Factory interface {
	Create(name string, options interface{}) (interrupt.Interrupt, error)
}
