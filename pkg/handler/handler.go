package handler

import "github.com/dpb587/slack-delegate-bot/pkg/message"

//go:generate counterfeiter . Handler
type Handler interface {
	Execute(*message.Message) (message.MessageResponse, error)
}
