package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	server *http.Server
}

func NewServer(server *http.Server) *Server {
	return &Server{
		server: server,
	}
}

func (s *Server) Run(services ...Service) error {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	for _, service := range services {
		service.InstallService(e)
	}

	return e.StartServer(s.server)
}
