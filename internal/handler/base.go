package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wzhanjun/go-echo-skeleton/internal/dto"
	"github.com/wzhanjun/go-echo-skeleton/internal/enum"
	slog "github.com/wzhanjun/log-service/client"
)

type BaseController struct {
}

// 请求成功
func (s BaseController) Success(c echo.Context, data interface{}, msg ...string) error {
	resp := dto.Response{
		Code: int32(enum.Success),
		Msg:  enum.Success.String(),
		Data: data,
	}
	if len(msg) > 0 {
		resp.Msg = msg[0]
	}
	return c.JSON(http.StatusOK, resp)
}

// 请求失败
func (s BaseController) Fail(c echo.Context, err error) error {
	slog.Label("request").Errorf("request failed. path:%s, error:%+v", c.Request().URL.Path, err)
	switch e := err.(type) {
	case enum.ErrCode:
		// 未知错误
		return c.JSON(http.StatusBadRequest, dto.Response{
			Code: int32(e),
			Msg:  e.String(),
		})
	case enum.ApiError:
		// 业务错误
		return c.JSON(http.StatusBadRequest, dto.Response{
			Code: int32(e.Code),
			Msg:  e.Msg,
		})
	case *echo.HTTPError:
		// echo 绑定参数错误
		msg := enum.ParamsError.String()
		if params, ok := err.(*echo.HTTPError).Internal.(*json.UnmarshalTypeError); ok {
			msg = fmt.Sprintf("params %s error", params.Field)
		}
		return c.JSON(http.StatusBadRequest, dto.Response{
			Code: int32(enum.ParamsError),
			Msg:  msg,
		})
	}

	return c.JSON(http.StatusBadRequest, dto.Response{
		Code: int32(enum.Error),
		Msg:  enum.Error.String(),
	})
}
