package literalmap

import (
	"github.com/dpb587/slack-delegate-bot/logic/delegate"
	"github.com/dpb587/slack-delegate-bot/delegatebot/message"
)

type Delegator struct {
	From       delegate.Delegator
	Users      map[string]string
	Usergroups map[string]string
}

var _ delegate.Delegator = &Delegator{}

func (i Delegator) Delegate(m message.Message) ([]delegate.Delegate, error) {
	inner, err := i.From.Delegate(m)
	if err != nil {
		return nil, err
	}

	var res []delegate.Delegate

	for _, innerInterrupt := range inner {
		literalInterrupt, ok := innerInterrupt.(delegate.Literal)
		if !ok {
			res = append(res, innerInterrupt)

			continue
		}

		var newres delegate.Delegate = literalInterrupt

		if mapped, found := i.Users[literalInterrupt.Text]; found {
			newres = delegate.User{ID: mapped}
		} else if mapped, found := i.Usergroups[literalInterrupt.Text]; found {
			newres = delegate.UserGroup{ID: mapped}
		}

		res = append(res, newres)
	}

	return res, nil
}
