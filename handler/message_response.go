package handler

import (
	"github.com/dpb587/slack-delegate-bot/logic/interrupt"
)

type MessageResponse struct {
	Interrupts   []interrupt.Interruptible
	EmptyMessage string
}
