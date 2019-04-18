package pairist

import (
	"github.com/dpb587/go-pairist/api"
	"github.com/dpb587/go-pairist/denormalized"
	"github.com/dpb587/slack-delegate-bot/interrupt"
	"github.com/dpb587/slack-delegate-bot/message"
)

type Interrupt struct {
	Team   string
	Role   string
	People map[string]string
}

var _ interrupt.Interrupt = &Interrupt{}

func (i Interrupt) Lookup(_ message.Message) ([]interrupt.Interruptible, error) {
	curr, err := api.DefaultClient.GetTeamCurrent(i.Team)
	if err != nil {
		return nil, err
	}

	var res []interrupt.Interruptible

	for _, lane := range denormalized.BuildLanes(curr).ByRole(i.Role) {
		for _, person := range lane.People {
			if handle, ok := i.People[person.Name]; ok {
				res = append(res, interrupt.User{ID: handle})
			} else {
				res = append(res, interrupt.Literal{Text: person.Name})
			}
		}
	}

	return res, nil
}
