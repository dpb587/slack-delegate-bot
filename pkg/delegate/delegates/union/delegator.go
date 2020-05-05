package union

import (
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
)

type Delegator struct {
	Delegators []delegate.Delegator
}

var _ delegate.Delegator = &Delegator{}

func (i Delegator) Delegate(m message.Message) ([]message.Delegate, error) {
	var r []message.Delegate

	for _, sub := range i.Delegators {
		subr, err := sub.Delegate(m)
		if err != nil {
			return nil, err
		}

		r = append(r, subr...)
	}

	return r, nil
}
