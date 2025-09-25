package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/wzhanjun/go-echo-skeleton/internal/handler"
)

type indexController struct {
	handler.BaseController
}

func NewIndexController() *indexController {
	return &indexController{}
}
func (s *indexController) Index(c echo.Context) error {
	return s.Success(c, "hello world")
}

func (s *indexController) Health(c echo.Context) error {
	return s.Success(c, nil)
}
