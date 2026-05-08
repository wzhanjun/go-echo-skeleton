package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wzhanjun/go-echo-skeleton/internal/dto"
	"github.com/wzhanjun/go-echo-skeleton/internal/handler"
	"github.com/wzhanjun/go-echo-skeleton/internal/models"
	"github.com/wzhanjun/go-echo-skeleton/internal/services"
)

type userController struct {
	handler.BaseController
	svc *services.UserService
}

func NewUserController() *userController {
	return &userController{svc: services.NewUserService()}
}

// GET /api/v1/users/:id
func (c *userController) Get(ctx echo.Context) error {
	user, err := c.svc.Get(ctx.Param("id"))
	if err != nil {
		return c.Fail(ctx, err)
	}
	return c.Success(ctx, user)
}

// GET /api/v1/users
func (c *userController) List(ctx echo.Context) error {
	var params dto.PageParams
	if err := ctx.Bind(&params); err != nil {
		return c.Fail(ctx, err)
	}

	users, total, err := c.svc.List(params)
	if err != nil {
		return c.Fail(ctx, err)
	}
	return c.Success(ctx, dto.NewPageData(total, users))
}

// POST /api/v1/users
func (c *userController) Create(ctx echo.Context) error {
	var user models.User
	if err := ctx.Bind(&user); err != nil {
		return c.Fail(ctx, err)
	}
	if user.Username == "" {
		return ctx.JSON(http.StatusBadRequest, dto.Response{
			Code: 1001,
			Msg:  "username is required",
		})
	}
	if err := c.svc.Create(&user); err != nil {
		return c.Fail(ctx, err)
	}
	return c.Success(ctx, user)
}

// PUT /api/v1/users/:id
func (c *userController) Update(ctx echo.Context) error {
	var user models.User
	if err := ctx.Bind(&user); err != nil {
		return c.Fail(ctx, err)
	}
	if err := c.svc.Update(ctx.Param("id"), &user); err != nil {
		return c.Fail(ctx, err)
	}
	return c.Success(ctx, nil)
}

// DELETE /api/v1/users/:id
func (c *userController) Delete(ctx echo.Context) error {
	if err := c.svc.Delete(ctx.Param("id")); err != nil {
		return c.Fail(ctx, err)
	}
	return c.Success(ctx, nil)
}
