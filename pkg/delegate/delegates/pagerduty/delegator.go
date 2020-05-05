package pagerduty

import (
	pagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
	"github.com/pkg/errors"
)

type Delegator struct {
	Client           *pagerduty.Client
	EscalationPolicy string
	EscalationLevel  uint
}

var _ delegate.Delegator = &Delegator{}

func (i Delegator) Delegate(_ message.Message) ([]message.Delegate, error) {
	oncalls, err := i.Client.ListOnCalls(pagerduty.ListOnCallOptions{
		Includes:            []string{"users"},
		EscalationPolicyIDs: []string{i.EscalationPolicy},
	})
	if err != nil {
		return nil, errors.Wrap(err, "listing on-calls")
	}

	scheduledUsers := map[string]struct{}{}

	for _, oncall := range oncalls.OnCalls {
		if i.EscalationLevel > 0 && i.EscalationLevel != oncall.EscalationLevel {
			continue
		}

		scheduledUsers[oncall.User.Email] = struct{}{}
	}

	var res []message.Delegate

	for scheduledUser := range scheduledUsers {
		res = append(res, delegate.Literal{Text: scheduledUser})
	}

	return res, nil
}
