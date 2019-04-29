package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/config"
	"github.com/dpb587/slack-delegate-bot/logic/delegate"
	"github.com/dpb587/slack-delegate-bot/logic/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/logic/delegate/delegates/union"
	"github.com/pkg/errors"
)

type factory struct {
	factory delegates.Factory
}

type Options []interface{}

func New(ff delegates.Factory) delegates.Factory {
	return &factory{
		factory: ff,
	}
}

func (f factory) Create(name string, options interface{}) (delegate.Delegator, error) {
	if name != "union" {
		return nil, fmt.Errorf("invalid delegate: %s", name)
	}

	parsed := Options{}

	err := config.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	var ccds []delegate.Delegator

	for optionsIdx, options := range parsed {
		key, value, err := config.KeyValueTuple(options)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing union delegate %d", optionsIdx)
		}

		delegate, err := f.factory.Create(key, value)
		if err != nil {
			return nil, errors.Wrapf(err, "creating union delegate %d", optionsIdx)
		}

		ccds = append(ccds, delegate)
	}

	return &union.Delegator{
		Delegators: ccds,
	}, nil
}
