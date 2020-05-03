package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MetaRuntimeHandler struct{}

var _ Service = MetaRuntimeHandler{}

func (h MetaRuntimeHandler) InstallService(e *echo.Echo) {
	r := e.Group("/api/v1/meta/runtime")

	r.GET("/failure", h.Failure)
	r.GET("/success", h.Success)
}

func (h MetaRuntimeHandler) Failure(c echo.Context) error {
	return errors.New("failure")
}

func (h MetaRuntimeHandler) Success(c echo.Context) error {
	return c.String(http.StatusOK, "success")
}
