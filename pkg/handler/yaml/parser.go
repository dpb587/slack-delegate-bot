package yaml

import (
	"github.com/dpb587/slack-delegate-bot/pkg/condition/conditions"
	"github.com/dpb587/slack-delegate-bot/pkg/configutil"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/pkg/handler"
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

func (l Parser) Parse(buf []byte) (handler.Handler, error) {
	h := Handler{}

	var parsed Schema

	err := yaml.Unmarshal(buf, &parsed.DelegateBot) // TODO remove `delegatebot` wrapper?
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling")
	}

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
