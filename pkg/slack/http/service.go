package http

import (
	"github.com/dpb587/slack-delegate-bot/pkg/http"
	"github.com/dpb587/slack-delegate-bot/pkg/slack"
	"github.com/dpb587/slack-delegate-bot/pkg/slack/slash"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Service struct {
	EventProcessor slack.Processor
	SlashProcessor slash.Processor
	SigningSecret  string
	Logger         *zap.Logger
}

var _ http.Service = &Service{}

func (s *Service) InstallService(e *echo.Echo) {
	r := e.Group("/api/v1/slack")

	{
		r := r.Group("/event")
		h := NewEventHandler(s.EventProcessor, s.SigningSecret, s.Logger)

		r.POST("", h.Accept)
	}

	{
		r := r.Group("/slash")
		h := NewSlashHandler(s.SlashProcessor, s.SigningSecret, s.Logger)

		r.POST("", h.Accept)
	}
}
