package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wzhanjun/go-echo-skeleton/internal/dto"
	"github.com/wzhanjun/go-echo-skeleton/internal/enum"
	slog "github.com/wzhanjun/log-service/client"
)

func DefaultHTTPErrorHandler(err error, c echo.Context) {
	slog.Std().Errorf("system error:%s, path:%s", err, c.Request().RequestURI)

	errs := enum.Error
	switch {
	case errors.Is(err, echo.ErrNotFound):
		errs = enum.ErrCodeNotFound
	}

	_ = c.JSON(http.StatusBadRequest, dto.Response{
		Code: int32(errs),
		Msg:  errs.String(),
		Data: nil,
	})
}
