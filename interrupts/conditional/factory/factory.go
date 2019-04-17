package factory

import (
	"fmt"

	"github.com/dpb587/slack-alias-bot/conditions"
	"github.com/dpb587/slack-alias-bot/config"
	"github.com/dpb587/slack-alias-bot/interrupt"
	"github.com/dpb587/slack-alias-bot/interrupts"
	"github.com/dpb587/slack-alias-bot/interrupts/conditional"
	"github.com/pkg/errors"
)

type factory struct {
	interruptsFactory interrupts.Factory
	conditionsFactory conditions.Factory
}

type Options struct {
	When interface{}                 `yaml:"when"`
	Then map[interface{}]interface{} `yaml:"then"`
	Else map[interface{}]interface{} `yaml:"else"`
}

func New(interruptsFactory interrupts.Factory, conditionsFactory conditions.Factory) interrupts.Factory {
	return &factory{
		interruptsFactory: interruptsFactory,
		conditionsFactory: conditionsFactory,
	}
}

func (f factory) Create(name string, options interface{}) (interrupt.Interrupt, error) {
	if name != "if" {
		return nil, fmt.Errorf("unsupported interrupt: %s", name)
	}

	parsed := Options{}

	err := config.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	when, err := f.conditionsFactory.Create("and", parsed.When)
	if err != nil {
		return nil, errors.Wrap(err, "creating conditional when")
	}

	thenName, thenOptions, err := config.KeyValueTuple(parsed.Then)
	if err != nil {
		return nil, errors.Wrap(err, "parsing conditional then")
	}

	then, err := f.interruptsFactory.Create(thenName, thenOptions)
	if err != nil {
		return nil, errors.Wrap(err, "creating conditional then")
	}

	var else_ interrupt.Interrupt

	if parsed.Then != nil {
		elseName, elseOptions, err := config.KeyValueTuple(parsed.Then)
		if err != nil {
			return nil, errors.Wrap(err, "parsing conditional else")
		}

		else_, err = f.interruptsFactory.Create(elseName, elseOptions)
		if err != nil {
			return nil, errors.Wrap(err, "creating conditional else")
		}
	}

	return &conditional.Interrupt{
		When: when,
		Then: then,
		Else: else_,
	}, nil
}
