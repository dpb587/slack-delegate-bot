package http

import (
	"github.com/dpb587/slack-delegate-bot/pkg/http"
	"github.com/dpb587/slack-delegate-bot/pkg/slack"
	"github.com/labstack/echo/v4"
)

type Service struct {
	Processor     slack.Processor
	SigningSecret string
}

var _ http.Service = &Service{}

func (s *Service) InstallService(e *echo.Echo) {
	r := e.Group("/api/v1/slack")

	{
		r := r.Group("/event")
		h := NewEventHandler(s.Processor, s.SigningSecret)

		r.POST("", h.Accept)
	}
}
