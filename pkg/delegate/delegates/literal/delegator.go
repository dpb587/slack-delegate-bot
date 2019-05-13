package literal

import (
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/message"
)

type Delegator struct {
	Text string
}

var _ delegate.Delegator = &Delegator{}

func (i Delegator) Delegate(_ message.Message) ([]delegate.Delegate, error) {
	return []delegate.Delegate{delegate.Literal{Text: i.Text}}, nil
}
