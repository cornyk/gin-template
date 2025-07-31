package controllers

import (
	"cornyk/gin-template/internal/exceptions"
	"cornyk/gin-template/internal/services"
	"cornyk/gin-template/internal/utils/response_util"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (c *UserController) GetUsers(ctx *gin.Context) {
	userService := services.UserService{}

	users, err := userService.GetAllUsers(ctx)
	if err != nil {
		ctx.Error(&exceptions.BusinessError{Message: "test", Code: 123})
		return
	}
	response_util.SucJson(ctx, response_util.Pagination{
		List:  users,
		Count: len(users),
	})
}
