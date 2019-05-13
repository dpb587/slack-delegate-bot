package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/pkg/config"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/topiclookup"
	"github.com/pkg/errors"
)

type factory struct {
	api topiclookup.SlackAPI
}

type Options struct {
	Channel string `yaml:"channel"`
}

func New(api topiclookup.SlackAPI) delegates.Factory {
	return &factory{
		api: api,
	}
}

func (f factory) Create(name string, options interface{}) (delegate.Delegator, error) {
	if name != "topiclookup" {
		return nil, fmt.Errorf("invalid delegate: %s", name)
	}

	parsed := Options{}

	err := config.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	return &topiclookup.Delegator{
		API:     f.api,
		Channel: parsed.Channel,
	}, nil
}
