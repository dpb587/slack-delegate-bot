package emaillookupmap

import (
	"strings"

	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/message"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/slack-go/slack"
)

type SlackAPI interface {
	GetUserByEmail(email string) (*slack.User, error)
}

type Delegator struct {
	From delegate.Delegator
	API  SlackAPI
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
		} else if !strings.Contains(literalInterrupt.Text, "@") {
			res = append(res, innerInterrupt)

			continue
		}

		user, err := i.API.GetUserByEmail(literalInterrupt.Text)
		if err != nil {
			// TODO warn?
		}

		if user == nil {
			res = append(res, innerInterrupt)

			continue
		}

		res = append(res, delegate.User{ID: user.ID})
	}

	return res, nil
}
