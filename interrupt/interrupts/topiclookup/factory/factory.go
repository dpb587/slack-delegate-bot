package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/config"
	"github.com/dpb587/slack-delegate-bot/interrupt"
	"github.com/dpb587/slack-delegate-bot/interrupt/interrupts"
	"github.com/dpb587/slack-delegate-bot/interrupt/interrupts/topiclookup"
	"github.com/pkg/errors"
)

type factory struct {
	api topiclookup.SlackAPI
}

type Options struct {
	Channel string `yaml:"channel"`
}

func New(api topiclookup.SlackAPI) interrupts.Factory {
	return &factory{
		api: api,
	}
}

func (f factory) Create(name string, options interface{}) (interrupt.Interrupt, error) {
	if name != "topiclookup" {
		return nil, fmt.Errorf("invalid interrupt: %s", name)
	}

	parsed := Options{}

	err := config.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	return &topiclookup.Interrupt{
		API:     f.api,
		Channel: parsed.Channel,
	}, nil
}
