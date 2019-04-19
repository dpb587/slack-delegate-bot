package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/config"
	"github.com/dpb587/slack-delegate-bot/logic/interrupt"
	"github.com/dpb587/slack-delegate-bot/logic/interrupt/interrupts"
	"github.com/dpb587/slack-delegate-bot/logic/interrupt/interrupts/literalmap"
	"github.com/pkg/errors"
)

type factory struct {
	interruptsFactory interrupts.Factory
}

type Options struct {
	From       map[interface{}]interface{} `yaml:"from"`
	Users      map[string]string           `yaml:"users"`
	Usergroups map[string]string           `yaml:"usergroups"`
}

func New(interruptsFactory interrupts.Factory) interrupts.Factory {
	return &factory{
		interruptsFactory: interruptsFactory,
	}
}

func (f factory) Create(name string, options interface{}) (interrupt.Interrupt, error) {
	if name != "literalmap" {
		return nil, fmt.Errorf("unsupported interrupt: %s", name)
	}

	parsed := Options{}

	err := config.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	fromName, fromOptions, err := config.KeyValueTuple(parsed.From)

	from, err := f.interruptsFactory.Create(fromName, fromOptions)
	if err != nil {
		return nil, errors.Wrap(err, "creating literalmap from")
	}

	return &literalmap.Interrupt{
		From:       from,
		Users:      parsed.Users,
		Usergroups: parsed.Usergroups,
	}, nil
}
