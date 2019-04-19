package handler

import (
	"github.com/dpb587/slack-delegate-bot/logic/delegate"
)

type MessageResponse struct {
	Delegates    []delegate.Delegate
	EmptyMessage string
}
