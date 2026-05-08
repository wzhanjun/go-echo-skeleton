package routers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	v1 "github.com/wzhanjun/go-echo-skeleton/internal/handler/api/v1"
	mid "github.com/wzhanjun/go-echo-skeleton/internal/middleware"
)

func Router() *echo.Echo {

	e := echo.New()

	// 内部http错误处理
	e.HTTPErrorHandler = mid.DefaultHTTPErrorHandler
	// 错误恢复
	e.Use(middleware.Recover(), middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
		MaxAge:           3600,
	}))

	e.GET("/health", v1.NewIndexController().Health)

	user := v1.NewUserController()
	apiv1 := e.Group("/api/v1")
	apiv1.Use(mid.RequestTime())
	{
		apiv1.GET("/index", v1.NewIndexController().Index)
		apiv1.GET("/users", user.List)
		apiv1.GET("/users/:id", user.Get)
		apiv1.POST("/users", user.Create)
		apiv1.PUT("/users/:id", user.Update)
		apiv1.DELETE("/users/:id", user.Delete)
	}

	return e
}
