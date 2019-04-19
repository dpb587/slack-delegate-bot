package topiclookup

import (
	"regexp"
	"strings"

	"github.com/dpb587/slack-delegate-bot/logic/interrupt"
	"github.com/dpb587/slack-delegate-bot/message"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

//go:generate counterfeiter . SlackAPI
type SlackAPI interface {
	GetChannelInfo(string) (*slack.Channel, error)
}

type Interrupt struct {
	API     SlackAPI
	Channel string
}

var _ interrupt.Interrupt = &Interrupt{}

var slackRefRE = regexp.MustCompile(`<[^>]+>`)
var topicInterruptREs = []*regexp.Regexp{
	regexp.MustCompile(`(?i)[*_]*interrupt[*_:]* (<[^>]+>\s*)+`),
}

func (i Interrupt) Lookup(m message.Message) ([]interrupt.Interruptible, error) {
	channel := m.InterruptTarget
	if i.Channel != "" {
		channel = i.Channel
	}

	info, err := i.API.GetChannelInfo(channel)
	if err != nil {
		return nil, errors.Wrapf(err, "getting info of channel %s", channel)
	}

	for _, topicInterruptRE := range topicInterruptREs {
		matches := topicInterruptRE.FindStringSubmatch(info.Topic.Value)
		if len(matches) == 0 {
			continue
		}

		slackRefMatches := slackRefRE.FindAllStringSubmatch(matches[0], -1)

		var results []interrupt.Interruptible

		for _, slackRefMatch := range slackRefMatches {
			match := slackRefMatch[0]

			if strings.HasPrefix(match, "<!subteam^") {
				pieces := strings.SplitN(strings.TrimPrefix(strings.TrimSuffix(strings.TrimSpace(match), ">"), "<!subteam^"), "|", 2)
				if len(pieces) != 2 {
					continue
				}

				results = append(results, interrupt.UserGroup{ID: pieces[0], Alias: strings.TrimPrefix(pieces[1], "@")})
			} else if strings.HasPrefix(match, "<@U") {
				results = append(results, interrupt.User{ID: strings.TrimPrefix(strings.TrimSuffix(strings.TrimSpace(match), ">"), "<@")})
			}
		}

		return results, nil
	}

	return nil, nil
}
