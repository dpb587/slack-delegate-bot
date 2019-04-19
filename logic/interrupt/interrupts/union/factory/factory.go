package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/config"
	"github.com/dpb587/slack-delegate-bot/logic/interrupt"
	"github.com/dpb587/slack-delegate-bot/logic/interrupt/interrupts"
	"github.com/dpb587/slack-delegate-bot/logic/interrupt/interrupts/union"
	"github.com/pkg/errors"
)

type factory struct {
	factory interrupts.Factory
}

type Options []interface{}

func New(ff interrupts.Factory) interrupts.Factory {
	return &factory{
		factory: ff,
	}
}

func (f factory) Create(name string, options interface{}) (interrupt.Interrupt, error) {
	if name != "union" {
		return nil, fmt.Errorf("invalid interrupt: %s", name)
	}

	parsed := Options{}

	err := config.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	var ccds []interrupt.Interrupt

	for optionsIdx, options := range parsed {
		key, value, err := config.KeyValueTuple(options)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing union interrupt %d", optionsIdx)
		}

		interrupt, err := f.factory.Create(key, value)
		if err != nil {
			return nil, errors.Wrapf(err, "creating union interrupt %d", optionsIdx)
		}

		ccds = append(ccds, interrupt)
	}

	return &union.Interrupt{
		Interrupts: ccds,
	}, nil
}
