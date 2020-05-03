package topiclookup

import (
	"regexp"
	"strings"

	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

//go:generate counterfeiter . SlackAPI
type SlackAPI interface {
	GetConversationInfo(string, bool) (*slack.Channel, error)
}

type Delegator struct {
	Channel string
}

var _ delegate.Delegator = &Delegator{}

var slackRefRE = regexp.MustCompile(`<[^>]+>`)
var topicInterruptREs = []*regexp.Regexp{
	regexp.MustCompile("(?i)[`*_]*interrupt[`*_:]*\\s+(<[^>]+>(,?\\s+and\\s+|,?\\s*)?)+"),
}

func (i Delegator) Delegate(m message.Message) ([]delegate.Delegate, error) {
	channel := m.InterruptTarget

	api, ok := m.ServiceAPI.(SlackAPI)
	if !ok {
		return nil, nil
	}

	info, err := api.GetConversationInfo(channel, false)
	if err != nil {
		return nil, errors.Wrapf(err, "getting info of channel %s", channel)
	}

	for _, topicInterruptRE := range topicInterruptREs {
		matches := topicInterruptRE.FindStringSubmatch(info.Topic.Value)
		if len(matches) == 0 {
			continue
		}

		slackRefMatches := slackRefRE.FindAllStringSubmatch(matches[0], -1)

		var results []delegate.Delegate

		for _, slackRefMatch := range slackRefMatches {
			match := slackRefMatch[0]

			if strings.HasPrefix(match, "<!subteam^") {
				pieces := strings.SplitN(strings.TrimPrefix(strings.TrimSuffix(strings.TrimSpace(match), ">"), "<!subteam^"), "|", 2)
				if len(pieces) != 2 {
					continue
				}

				results = append(results, delegate.UserGroup{ID: pieces[0], Alias: strings.TrimPrefix(pieces[1], "@")})
			} else if strings.HasPrefix(match, "<@U") {
				results = append(results, delegate.User{ID: strings.TrimPrefix(strings.TrimSuffix(strings.TrimSpace(match), ">"), "<@")})
			}
		}

		return results, nil
	}

	return nil, nil
}
