package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/pkg/condition"
	"github.com/dpb587/slack-delegate-bot/pkg/condition/conditions"
	"github.com/dpb587/slack-delegate-bot/pkg/condition/conditions/boolnot"
	"github.com/dpb587/slack-delegate-bot/pkg/config"
	"github.com/pkg/errors"
)

type factory struct {
	factory conditions.Factory
}

func New(ff conditions.Factory) conditions.Factory {
	return &factory{
		factory: ff,
	}
}

func (f factory) Create(name string, options interface{}) (condition.Condition, error) {
	if name != "not" {
		return nil, fmt.Errorf("invalid condition: %s", name)
	}

	key, value, err := config.KeyValueTuple(options)
	if err != nil {
		return nil, errors.Wrap(err, "parsing")
	}

	condition, err := f.factory.Create(key, value)
	if err != nil {
		return nil, errors.Wrap(err, "creating")
	}

	return &boolnot.Condition{
		Condition: condition,
	}, nil
}
