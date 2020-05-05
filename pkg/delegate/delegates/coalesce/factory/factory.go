package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/pkg/configutil"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/coalesce"
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
	if name != "coalesce" {
		return nil, fmt.Errorf("invalid delegate: %s", name)
	}

	parsed := Options{}

	err := configutil.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	var ccds []delegate.Delegator

	for optionsIdx, options := range parsed {
		key, value, err := configutil.KeyValueTuple(options)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing coalesce delegate %d", optionsIdx)
		}

		delegate, err := f.factory.Create(key, value)
		if err != nil {
			return nil, errors.Wrapf(err, "creating coalesce delegate %d", optionsIdx)
		}

		ccds = append(ccds, delegate)
	}

	return &coalesce.Delegator{
		Delegators: ccds,
	}, nil
}
