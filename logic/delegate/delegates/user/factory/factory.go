package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/config"
	"github.com/dpb587/slack-delegate-bot/logic/delegate"
	"github.com/dpb587/slack-delegate-bot/logic/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/logic/delegate/delegates/user"
	"github.com/pkg/errors"
)

type factory struct{}

type Options struct {
	ID string `yaml:"id"`
}

func New() delegates.Factory {
	return &factory{}
}

func (f factory) Create(name string, options interface{}) (delegate.Delegator, error) {
	if name != "user" {
		return nil, fmt.Errorf("unsupported delegate: %s", name)
	}

	parsed := Options{}

	err := config.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	return &user.Delegator{
		ID: parsed.ID,
	}, nil
}
