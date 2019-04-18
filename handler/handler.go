package handler

import "github.com/dpb587/slack-delegate-bot/message"

type Handler interface {
	IsApplicable(message.Message) (bool, error)
	Apply(*message.Message) error
}
