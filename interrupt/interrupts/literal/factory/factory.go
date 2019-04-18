package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/config"
	"github.com/dpb587/slack-delegate-bot/interrupt"
	"github.com/dpb587/slack-delegate-bot/interrupt/interrupts"
	"github.com/dpb587/slack-delegate-bot/interrupt/interrupts/literal"
	"github.com/pkg/errors"
)

type factory struct{}

type Options struct {
	Text string `yaml:"text"`
}

func New() interrupts.Factory {
	return &factory{}
}

func (f factory) Create(name string, options interface{}) (interrupt.Interrupt, error) {
	if name != "literal" {
		return nil, fmt.Errorf("unsupported interrupt: %s", name)
	}

	parsed := Options{}

	err := config.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	return &literal.Interrupt{
		Text: parsed.Text,
	}, nil
}
