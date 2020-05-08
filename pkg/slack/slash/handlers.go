package slash

import (
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type Handlers []Handler

var _ Handler = Handlers{}

func (hh Handlers) UsageHint() string {
	return ""
}

func (hh Handlers) ShortDescription() string {
	return ""
}

func (hh Handlers) Handle(cmd slack.SlashCommand) (bool, error) {
	for hIdx, h := range hh {
		done, err := h.Handle(cmd)
		if err != nil {
			return false, errors.Wrapf(err, "processing handler %d", hIdx)
		} else if done {
			return true, nil
		}
	}

	return false, nil
}
