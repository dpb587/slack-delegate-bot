package http

import (
	"github.com/dpb587/slack-delegate-bot/pkg/http"
	"github.com/dpb587/slack-delegate-bot/pkg/slack"
	"github.com/dpb587/slack-delegate-bot/pkg/slack/slash"
	"github.com/labstack/echo/v4"
)

type Service struct {
	EventProcessor slack.Processor
	SlashProcessor slash.Processor
	SigningSecret  string
}

var _ http.Service = &Service{}

func (s *Service) InstallService(e *echo.Echo) {
	r := e.Group("/api/v1/slack")

	{
		r := r.Group("/event")
		h := NewEventHandler(s.EventProcessor, s.SigningSecret)

		r.POST("", h.Accept)
	}

	{
		r := r.Group("/slash")
		h := NewSlashHandler(s.SlashProcessor, s.SigningSecret)

		r.POST("", h.Accept)
	}
}
