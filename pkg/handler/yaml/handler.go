package yaml

import (
	"github.com/dpb587/slack-delegate-bot/pkg/condition"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/handler"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
	"github.com/pkg/errors"
)

type Handler struct {
	condition condition.Condition
	delegator delegate.Delegator
	options   SchemaDelegateBotWithOptions
}

var _ handler.Handler = &Handler{}

func (h Handler) Execute(m *message.Message) (message.MessageResponse, error) {
	if h.condition != nil {
		tf, err := h.condition.Evaluate(*m)
		if err != nil {
			return message.MessageResponse{}, errors.Wrap(err, "evaluating condition")
		} else if !tf {
			return message.MessageResponse{}, nil
		}
	}

	response := message.MessageResponse{}

	interrupts, err := h.delegator.Delegate(*m)
	if err != nil {
		return response, err
	}

	if len(interrupts) == 0 {
		response.EmptyMessage = h.options.EmptyMessage
	} else {
		response.Delegates = interrupts
	}

	return response, nil
}
