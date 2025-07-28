package controllers

import (
	"cornyk/gin-template/internal/daos"
	"cornyk/gin-template/internal/exceptions"
	"cornyk/gin-template/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// GetUsers 获取所有用户
func GetUsers(c *gin.Context) {
	users, err := daos.GetAllUsers(c)
	if err != nil {
		c.Error(&exceptions.BusinessError{Message: "test", Code: 123})
		return
	}
	response.SucJson(c, response.Pagination{
		List:  users,
		Count: len(users),
	})
}
