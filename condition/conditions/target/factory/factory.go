package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/condition"
	"github.com/dpb587/slack-delegate-bot/condition/conditions"
	"github.com/dpb587/slack-delegate-bot/condition/conditions/target"
	"github.com/dpb587/slack-delegate-bot/config"
	"github.com/pkg/errors"
)

type factory struct{}

type Options struct {
	Channel string `yaml:"channel"`
}

func New() conditions.Factory {
	return &factory{}
}

func (f factory) Create(name string, options interface{}) (condition.Condition, error) {
	if name != "target" {
		return nil, fmt.Errorf("invalid condition: %s", name)
	}

	parsed := Options{}

	err := config.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	return &target.Condition{
		Channel: parsed.Channel,
	}, nil
}
