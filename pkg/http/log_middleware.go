package http

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func LogMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()

			begin := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}

			end := time.Now()

			fields := []zap.Field{
				zap.String("remote_addr", req.RemoteAddr),
				zap.String("host", req.Host),
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.String("user_agent", req.UserAgent()),
				zap.String("referer", req.Referer()),
				zap.Int("status", res.Status),
				zap.Int64("duration", int64(end.Sub(begin)/ time.Millisecond)),
				zap.Int64("bytes_out", res.Size),
			}

			if err != nil {
				fields = append(fields, zap.Error(err))
			}

			logger.Info("http request finished", fields...)

			return
		}
	}
}