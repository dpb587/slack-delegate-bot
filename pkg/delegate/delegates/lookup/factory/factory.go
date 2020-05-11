package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/pkg/configutil"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/lookup"
	"github.com/pkg/errors"
)

type factory struct{}

type Options struct {
	Channel string `yaml:"channel"`
}

func New() delegates.Factory {
	return factory{}
}

func (f factory) Create(name string, options interface{}) (delegate.Delegator, error) {
	if name != "lookup" {
		return nil, fmt.Errorf("invalid delegate: %s", name)
	}

	parsed := Options{}

	err := configutil.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	return &lookup.Delegator{
		Channel: parsed.Channel,
	}, nil
}
