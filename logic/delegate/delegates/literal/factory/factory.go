package factory

import (
	"fmt"

	"github.com/dpb587/slack-delegate-bot/config"
	"github.com/dpb587/slack-delegate-bot/logic/delegate"
	"github.com/dpb587/slack-delegate-bot/logic/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/logic/delegate/delegates/literal"
	"github.com/pkg/errors"
)

type factory struct{}

type Options struct {
	Text string `yaml:"text"`
}

func New() delegates.Factory {
	return &factory{}
}

func (f factory) Create(name string, options interface{}) (delegate.Delegator, error) {
	if name != "literal" {
		return nil, fmt.Errorf("unsupported delegate: %s", name)
	}

	parsed := Options{}

	err := config.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	return &literal.Delegator{
		Text: parsed.Text,
	}, nil
}
