package pairist

import (
	"github.com/dpb587/go-pairist/api"
	"github.com/dpb587/go-pairist/denormalized"
	"github.com/dpb587/slack-delegate-bot/logic/delegate"
	"github.com/dpb587/slack-delegate-bot/delegatebot/message"
)

type Delegator struct {
	Team string
	Role string
}

var _ delegate.Delegator = &Delegator{}

func (i Delegator) Delegate(_ message.Message) ([]delegate.Delegate, error) {
	curr, err := api.DefaultClient.GetTeamCurrent(i.Team)
	if err != nil {
		return nil, err
	}

	var res []delegate.Delegate

	for _, lane := range denormalized.BuildLanes(curr).ByRole(i.Role) {
		for _, person := range lane.People {
			res = append(res, delegate.Literal{Text: person.Name})
		}
	}

	return res, nil
}
