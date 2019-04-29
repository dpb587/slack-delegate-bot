package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/config"
	"github.com/dpb587/slack-delegate-bot/logic/delegate"
	"github.com/dpb587/slack-delegate-bot/logic/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/logic/delegate/delegates/literalmap"
	"github.com/pkg/errors"
)

type factory struct {
	delegatesFactory delegates.Factory
}

type Options struct {
	From       map[interface{}]interface{} `yaml:"from"`
	Users      map[string]string           `yaml:"users"`
	Usergroups map[string]string           `yaml:"usergroups"`
}

func New(delegatesFactory delegates.Factory) delegates.Factory {
	return &factory{
		delegatesFactory: delegatesFactory,
	}
}

func (f factory) Create(name string, options interface{}) (delegate.Delegator, error) {
	if name != "literalmap" {
		return nil, fmt.Errorf("unsupported delegate: %s", name)
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

	return &literalmap.Delegator{
		From:       from,
		Users:      parsed.Users,
		Usergroups: parsed.Usergroups,
	}, nil
}
