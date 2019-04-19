package single

import (
	"github.com/dpb587/slack-delegate-bot/logic/condition"
	"github.com/dpb587/slack-delegate-bot/handler"
	"github.com/dpb587/slack-delegate-bot/logic/interrupt"
	"github.com/dpb587/slack-delegate-bot/message"
)

type Handler struct {
	Condition condition.Condition
	Interrupt interrupt.Interrupt
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

	interrupts, err := h.Interrupt.Lookup(*m)
	if err != nil {
		return response, err
	}

	if len(interrupts) == 0 {
		response.EmptyMessage = h.Options.EmptyMessage
	} else {
		response.Interrupts = interrupts
	}

	return response, nil
}
