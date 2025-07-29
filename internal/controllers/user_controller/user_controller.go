package user_controller

import (
	"cornyk/gin-template/internal/daos/user_dao"
	"cornyk/gin-template/internal/exceptions"
	"cornyk/gin-template/internal/utils/response_util"

	"github.com/gin-gonic/gin"
)

// GetUsers 获取所有用户
func GetUsers(c *gin.Context) {
	users, err := user_dao.GetAllUsers(c)
	if err != nil {
		c.Error(&exceptions.BusinessError{Message: "test", Code: 123})
		return
	}
	response_util.SucJson(c, response_util.Pagination{
		List:  users,
		Count: len(users),
	})
}
