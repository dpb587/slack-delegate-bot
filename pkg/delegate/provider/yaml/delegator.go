package yaml

import (
	"github.com/dpb587/slack-delegate-bot/pkg/condition"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
	"github.com/pkg/errors"
)

type Delegator struct {
	condition condition.Condition
	delegator delegate.Delegator
	options   SchemaDelegateBotWithOptions
}

var _ delegate.Delegator = &Delegator{}

func (h Delegator) Delegate(m message.Message) ([]message.Delegate, error) {
	if h.condition != nil {
		tf, err := h.condition.Evaluate(m)
		if err != nil {
			return nil, errors.Wrap(err, "evaluating condition")
		} else if !tf {
			return nil, nil
		}
	}

	res, err := h.delegator.Delegate(m)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 && len(h.options.EmptyMessage) > 0 {
		res = append(
			res,
			delegate.Literal{
				Text: h.options.EmptyMessage,
			},
		)
	}

	return res, nil
}
