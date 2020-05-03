package single

import (
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/handler"
	"github.com/dpb587/slack-delegate-bot/pkg/condition"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
)

type Handler struct {
	Condition condition.Condition
	Delegator delegate.Delegator
	Options   handler.Options
}

var _ handler.Handler = &Handler{}

func (h Handler) IsApplicable(m message.Message) (bool, error) {
	if h.Condition != nil {
		return h.Condition.Evaluate(m)
	}

	return true, nil
}

func (h Handler) Execute(m *message.Message) (handler.MessageResponse, error) {
	response := handler.MessageResponse{}

	interrupts, err := h.Delegator.Delegate(*m)
	if err != nil {
		return response, err
	}

	if len(interrupts) == 0 {
		response.EmptyMessage = h.Options.EmptyMessage
	} else {
		response.Delegates = interrupts
	}

	return response, nil
}
