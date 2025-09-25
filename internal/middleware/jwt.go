package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/wzhanjun/go-echo-skeleton/internal/dto"
	"github.com/wzhanjun/go-echo-skeleton/internal/enum"
	"github.com/wzhanjun/go-echo-skeleton/pkg/config"
)

func JwtAuth() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.Cfg.JWT.Secret),
		Skipper: func(c echo.Context) bool {
			return false
		},
		ErrorHandler: func(c echo.Context, err error) error {
			_ = c.JSON(http.StatusUnauthorized, dto.Response{
				Code: int32(enum.ErrCodeUnauthorized),
				Msg:  enum.ErrCodeUnauthorized.String(),
			})
			return nil
		},
		SuccessHandler: func(context echo.Context) {
			token := context.Get("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)
			context.Set("id", claims["id"])
			context.Set("username", claims["username"])
		},
	})
}
