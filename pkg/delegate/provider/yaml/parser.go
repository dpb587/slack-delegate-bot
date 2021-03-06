package yaml

import (
	"github.com/dpb587/slack-delegate-bot/pkg/condition/conditions"
	"github.com/dpb587/slack-delegate-bot/pkg/configutil"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type Parser struct {
	delegatorsFactory delegates.Factory
	conditionsFactory conditions.Factory
}

func NewParser(delegatorsFactory delegates.Factory, conditionsFactory conditions.Factory) *Parser {
	return &Parser{
		delegatorsFactory: delegatorsFactory,
		conditionsFactory: conditionsFactory,
	}
}

func (l Parser) Parse(buf []byte) (delegate.Delegator, error) {
	var parsed Schema

	err := yaml.Unmarshal(buf, &parsed.DelegateBot)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling")
	}

	return l.parse(parsed)
}

func (l Parser) ParseFull(buf []byte) (delegate.Delegator, error) {
	var parsed Schema

	err := yaml.Unmarshal(buf, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling")
	}

	return l.parse(parsed)
}

func (l Parser) parse(parsed Schema) (delegate.Delegator, error) {
	h := Delegator{}

	if parsed.DelegateBot.Watch != nil {
		watch, err := l.conditionsFactory.Create("or", parsed.DelegateBot.Watch)
		if err != nil {
			return nil, errors.Wrap(err, "building watch")
		}

		h.condition = watch
	}

	delegateKey, delegateOptions, err := configutil.KeyValueTuple(parsed.DelegateBot.Delegate)
	if err != nil {
		return nil, errors.Wrap(err, "parsing delegate")
	}

	delegate, err := l.delegatorsFactory.Create(delegateKey, delegateOptions)
	if err != nil {
		return nil, errors.Wrap(err, "building delegate")
	}

	h.delegator = delegate

	h.options = parsed.DelegateBot.Options

	return h, nil
}
