package single

import (
	"fmt"

	"github.com/dpb587/slack-alias-bot/condition"
	"github.com/dpb587/slack-alias-bot/handler"
	"github.com/dpb587/slack-alias-bot/interrupt"
	"github.com/dpb587/slack-alias-bot/interrupts"
	"github.com/dpb587/slack-alias-bot/message"
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

func (h Handler) Apply(m *message.Message) error {
	lookups, err := h.Interrupt.Lookup(*m)
	if err != nil {
		return err
	}

	if len(lookups) == 0 {
		if h.Options.EmptyResponse != "" {
			m.SetResponse(message.Response{Text: h.Options.EmptyResponse})
		}

		return nil
	}

	m.SetResponse(message.Response{Text: fmt.Sprintf("^ %s", interrupts.Join(lookups, " "))})

	return nil
}
