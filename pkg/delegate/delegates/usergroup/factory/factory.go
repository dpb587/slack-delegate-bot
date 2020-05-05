package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/pkg/configutil"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/usergroup"
	"github.com/pkg/errors"
)

type factory struct{}

type Options struct {
	ID    string `yaml:"id"`
	Alias string `yaml:"alias"`
}

func New() delegates.Factory {
	return &factory{}
}

func (f factory) Create(name string, options interface{}) (delegate.Delegator, error) {
	if name != "usergroup" {
		return nil, fmt.Errorf("unsupported delegate: %s", name)
	}

	parsed := Options{}

	err := configutil.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	return &usergroup.Delegator{
		ID:    parsed.ID,
		Alias: parsed.Alias,
	}, nil
}
