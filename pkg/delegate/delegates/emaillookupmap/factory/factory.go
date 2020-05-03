package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/pkg/config"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/emaillookupmap"
	"github.com/pkg/errors"
)

type factory struct {
	delegatesFactory delegates.Factory
}

type Options struct {
	From map[interface{}]interface{} `yaml:"from"`
}

func New(delegatesFactory delegates.Factory) delegates.Factory {
	return &factory{
		delegatesFactory: delegatesFactory,
	}
}

func (f factory) Create(name string, options interface{}) (delegate.Delegator, error) {
	if name != "emaillookupmap" {
		return nil, fmt.Errorf("invalid delegate: %s", name)
	}

	parsed := Options{}

	err := config.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	fromName, fromOptions, err := config.KeyValueTuple(parsed.From)

	from, err := f.delegatesFactory.Create(fromName, fromOptions)
	if err != nil {
		return nil, errors.Wrap(err, "creating literalmap from")
	}

	return &emaillookupmap.Delegator{
		From: from,
	}, nil
}
