package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wzhanjun/go-echo-skeleton/internal/enum"
	slog "github.com/wzhanjun/log-service/client"
)

func RequestTime() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			if err := next(c); err != nil {
				c.Error(err)
			}
			slog.Label("request").Infof("[%s] | %d | %s | %s | %s | %s", time.Now().Format(enum.DateTimeTpl), c.Response().Status,
				time.Since(start), c.RealIP(), c.Request().Method, c.Request().URL)

			return nil
		}
	}
}
