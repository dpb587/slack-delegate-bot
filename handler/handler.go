package handler

import "github.com/dpb587/slack-alias-bot/message"

type Handler interface {
	IsApplicable(message.Message) (bool, error)
	Apply(*message.Message) error
}
