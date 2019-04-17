package multiple

import (
	"github.com/dpb587/slack-alias-bot/handler"
	"github.com/dpb587/slack-alias-bot/message"
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

func (h Handler) Apply(m *message.Message) error {
	// last one wins
	for _, hh := range h.Handlers {
		b, err := hh.IsApplicable(*m)
		if err != nil {
			return err
		} else if !b {
			continue
		}

		err = hh.Apply(m)
		if err != nil {
			return err
		}
	}

	return nil
}
