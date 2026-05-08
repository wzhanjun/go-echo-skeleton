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

	var (
		errs = enum.Error
		httpStatus = http.StatusInternalServerError
	)
	switch {
	case errors.Is(err, echo.ErrNotFound):
		errs = enum.ErrCodeNotFound
		httpStatus = http.StatusNotFound
	case errors.Is(err, echo.ErrUnauthorized):
		errs = enum.ErrCodeUnauthorized
		httpStatus = http.StatusUnauthorized
	case errors.Is(err, echo.ErrBadRequest):
		errs = enum.ParamsError
		httpStatus = http.StatusBadRequest
	}

	_ = c.JSON(httpStatus, dto.Response{
		Code: int32(errs),
		Msg:  errs.String(),
		Data: nil,
	})
}
