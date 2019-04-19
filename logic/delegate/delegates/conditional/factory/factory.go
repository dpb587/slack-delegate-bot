package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/config"
	"github.com/dpb587/slack-delegate-bot/logic/condition/conditions"
	"github.com/dpb587/slack-delegate-bot/logic/delegate"
	"github.com/dpb587/slack-delegate-bot/logic/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/logic/delegate/delegates/conditional"
	"github.com/pkg/errors"
)

type factory struct {
	delegatesFactory  delegates.Factory
	conditionsFactory conditions.Factory
}

type Options struct {
	When []interface{}               `yaml:"when"`
	Then map[interface{}]interface{} `yaml:"then"`
	Else map[interface{}]interface{} `yaml:"else"`
}

func New(delegatesFactory delegates.Factory, conditionsFactory conditions.Factory) delegates.Factory {
	return &factory{
		delegatesFactory:  delegatesFactory,
		conditionsFactory: conditionsFactory,
	}
}

func (f factory) Create(name string, options interface{}) (delegate.Delegator, error) {
	if name != "if" {
		return nil, fmt.Errorf("unsupported delegate: %s", name)
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

	then, err := f.delegatesFactory.Create(thenName, thenOptions)
	if err != nil {
		return nil, errors.Wrap(err, "creating conditional then")
	}

	var else_ delegate.Delegator

	if parsed.Else != nil {
		elseName, elseOptions, err := config.KeyValueTuple(parsed.Else)
		if err != nil {
			return nil, errors.Wrap(err, "parsing conditional else")
		}

		else_, err = f.delegatesFactory.Create(elseName, elseOptions)
		if err != nil {
			return nil, errors.Wrap(err, "creating conditional else")
		}
	}

	return &conditional.Delegator{
		When: when,
		Then: then,
		Else: else_,
	}, nil
}
