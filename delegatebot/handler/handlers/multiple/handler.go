package multiple

import (
	"github.com/dpb587/slack-delegate-bot/delegatebot/handler"
	"github.com/dpb587/slack-delegate-bot/delegatebot/message"
)

type Handler struct {
	Handlers []handler.Handler
}

var _ handler.Handler = &Handler{}

func (h Handler) IsApplicable(m message.Message) (bool, error) {
	for _, hh := range h.Handlers {
		b, err := hh.IsApplicable(m)
		if err != nil {
			return false, err
		} else if b {
			return b, nil
		}
	}

	return false, nil
}

func (h Handler) Execute(m *message.Message) (handler.MessageResponse, error) {
	response := handler.MessageResponse{}

	// first one finding interrupt wins
	for _, hh := range h.Handlers {
		b, err := hh.IsApplicable(*m)
		if err != nil {
			return handler.MessageResponse{}, err
		} else if !b {
			continue
		}

		r, err := hh.Execute(m)
		if err != nil {
			return handler.MessageResponse{}, err
		}

		if len(r.Delegates) > 0 {
			response.Delegates = r.Delegates

			// if interrupts are found, return immediately
			return response, nil
		}

		if response.EmptyMessage == "" && r.EmptyMessage != "" {
			response.EmptyMessage = r.EmptyMessage
		}
	}

	return response, nil
}