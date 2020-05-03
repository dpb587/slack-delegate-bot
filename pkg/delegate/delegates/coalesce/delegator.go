package coalesce

import (
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
)

type Delegator struct {
	Delegators []delegate.Delegator
}

var _ delegate.Delegator = &Delegator{}

func (i Delegator) Delegate(m message.Message) ([]delegate.Delegate, error) {
	for _, sub := range i.Delegators {
		subr, err := sub.Delegate(m)
		if err != nil {
			return nil, err
		} else if len(subr) > 0 {
			return subr, nil
		}
	}

	return nil, nil
}
