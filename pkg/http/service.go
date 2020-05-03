package http

import "github.com/labstack/echo/v4"

type Service interface {
	InstallService(e *echo.Echo)
}
