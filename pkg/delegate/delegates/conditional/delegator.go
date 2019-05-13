package conditional

import (
	"github.com/dpb587/slack-delegate-bot/pkg/condition"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/message"
)

type Delegator struct {
	When condition.Condition
	Then delegate.Delegator
	Else delegate.Delegator
}

var _ delegate.Delegator = &Delegator{}

func (i Delegator) Delegate(m message.Message) ([]delegate.Delegate, error) {
	when, err := i.When.Evaluate(m)
	if err != nil {
		return nil, err
	}

	if when {
		if i.Then != nil {
			return i.Then.Delegate(m)
		}
	} else if i.Else != nil {
		return i.Else.Delegate(m)
	}

	return nil, nil
}
