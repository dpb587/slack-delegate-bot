package pairist

import (
	"github.com/dpb587/go-pairist/api"
	"github.com/dpb587/go-pairist/denormalized"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
)

type Delegator struct {
	Client *api.Client

	Team string

	Role  string
	Track string
}

var _ delegate.Delegator = &Delegator{}

func (i Delegator) Delegate(_ message.Message) ([]delegate.Delegate, error) {
	curr, err := i.Client.GetTeamCurrent(i.Team)
	if err != nil {
		return nil, err
	}

	allLanes := denormalized.BuildLanes(curr)
	var lanes denormalized.Lanes

	if i.Track != "" {
		lanes = allLanes.ByTrack(i.Track)
	} else {
		lanes = allLanes.ByRole(i.Role)
	}

	var res []delegate.Delegate

	for _, lane := range lanes {
		for _, person := range lane.People {
			res = append(res, delegate.Literal{Text: person.Name})
		}
	}

	return res, nil
}
