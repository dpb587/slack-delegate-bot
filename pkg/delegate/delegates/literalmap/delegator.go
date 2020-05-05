package literalmap

import (
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
)

type Delegator struct {
	From       delegate.Delegator
	Users      map[string]string
	Usergroups map[string]string
}

var _ delegate.Delegator = &Delegator{}

func (i Delegator) Delegate(m message.Message) ([]message.Delegate, error) {
	inner, err := i.From.Delegate(m)
	if err != nil {
		return nil, err
	}

	var res []message.Delegate

	for _, innerInterrupt := range inner {
		literalInterrupt, ok := innerInterrupt.(delegate.Literal)
		if !ok {
			res = append(res, innerInterrupt)

			continue
		}

		var newres message.Delegate = literalInterrupt

		if mapped, found := i.Users[literalInterrupt.Text]; found {
			newres = delegate.User{ID: mapped}
		} else if mapped, found := i.Usergroups[literalInterrupt.Text]; found {
			newres = delegate.UserGroup{ID: mapped}
		}

		res = append(res, newres)
	}

	return res, nil
}
