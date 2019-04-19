package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/config"
	"github.com/dpb587/slack-delegate-bot/logic/interrupt"
	"github.com/dpb587/slack-delegate-bot/logic/interrupt/interrupts"
	"github.com/dpb587/slack-delegate-bot/logic/interrupt/interrupts/pairist"
	"github.com/pkg/errors"
)

type factory struct{}

type Options struct {
	Team string `yaml:"team"`
	Role string `yaml:"role"`
}

func New() interrupts.Factory {
	return &factory{}
}

func (f factory) Create(name string, options interface{}) (interrupt.Interrupt, error) {
	if name != "pairist" {
		return nil, fmt.Errorf("unsupported interrupt: %s", name)
	}

	parsed := Options{}

	err := config.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	return &pairist.Interrupt{
		Team: parsed.Team,
		Role: parsed.Role,
	}, nil
}
