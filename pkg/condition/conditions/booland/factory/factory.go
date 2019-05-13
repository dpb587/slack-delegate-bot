package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/pkg/condition"
	"github.com/dpb587/slack-delegate-bot/pkg/condition/conditions"
	"github.com/dpb587/slack-delegate-bot/pkg/condition/conditions/booland"
	"github.com/dpb587/slack-delegate-bot/pkg/config"
	"github.com/pkg/errors"
)

type factory struct {
	factory conditions.Factory
}

type Options []interface{}

func New(ff conditions.Factory) conditions.Factory {
	return &factory{
		factory: ff,
	}
}

func (f factory) Create(name string, options interface{}) (condition.Condition, error) {
	if name != "and" {
		return nil, fmt.Errorf("invalid condition: %s", name)
	}

	parsed := Options{}

	err := config.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	var ccds []condition.Condition

	for optionsIdx, options := range parsed {
		key, value, err := config.KeyValueTuple(options)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing condition %d", optionsIdx)
		}

		condition, err := f.factory.Create(key, value)
		if err != nil {
			return nil, errors.Wrapf(err, "creating condition %d", optionsIdx)
		}

		ccds = append(ccds, condition)
	}

	return &booland.Condition{
		Conditions: ccds,
	}, nil
}
