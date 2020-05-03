package multiple

import (
	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/handler"
	"github.com/dpb587/slack-delegate-bot/pkg/message"
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

	// first one matching wins
	for _, hh := range h.Handlers {
		b, err := hh.IsApplicable(*m)
		if err != nil {
			return handler.MessageResponse{}, err
		} else if !b {
			continue
		}

		return hh.Execute(m)
	}

	return response, nil
}
