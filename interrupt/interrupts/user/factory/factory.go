package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/config"
	"github.com/dpb587/slack-delegate-bot/interrupt"
	"github.com/dpb587/slack-delegate-bot/interrupt/interrupts"
	"github.com/dpb587/slack-delegate-bot/interrupt/interrupts/user"
	"github.com/pkg/errors"
)

type factory struct{}

type Options struct {
	ID string `yaml:"id"`
}

func New() interrupts.Factory {
	return &factory{}
}

func (f factory) Create(name string, options interface{}) (interrupt.Interrupt, error) {
	if name != "user" {
		return nil, fmt.Errorf("unsupported interrupt: %s", name)
	}

	parsed := Options{}

	err := config.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	return &user.Interrupt{
		ID: parsed.ID,
	}, nil
}
