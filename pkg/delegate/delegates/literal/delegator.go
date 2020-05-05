package literal

import (
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
)

type Delegator struct {
	Text string
}

var _ delegate.Delegator = &Delegator{}

func (i Delegator) Delegate(_ message.Message) ([]message.Delegate, error) {
	return []message.Delegate{delegate.Literal{Text: i.Text}}, nil
}
