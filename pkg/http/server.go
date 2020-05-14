package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
	logger *zap.Logger
}

func NewServer(server *http.Server, logger *zap.Logger) *Server {
	return &Server{
		server: server,
		logger: logger,
	}
}

func (s *Server) Run(services ...Service) error {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(LogMiddleware(s.logger))
	e.Use(middleware.Recover())

	for _, service := range services {
		service.InstallService(e)
	}

	return e.StartServer(s.server)
}
