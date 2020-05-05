package handler

import (
	"github.com/dpb587/slack-delegate-bot/pkg/message"
	"github.com/pkg/errors"
)

type CoalesceHandler struct {
	handlers []Handler
}

var _ Handler = &CoalesceHandler{}

func NewCoalesceHandler(handlers ...Handler) Handler {
	return &CoalesceHandler{
		handlers: handlers,
	}
}

func (ch *CoalesceHandler) Execute(msg *message.Message) (message.MessageResponse, error) {
	for hIdx, h := range ch.handlers {
		res, err := h.Execute(msg)
		if err != nil {
			// TODO optional ignore errors?
			return message.MessageResponse{}, errors.Wrapf(err, "executing handler %d", hIdx)
		} else if res.IsUnset() {
			continue
		}

		return res, nil
	}

	return message.MessageResponse{}, nil
}
