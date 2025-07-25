package controllers

import (
	"cornyk/gin-template/internal/daos"
	"cornyk/gin-template/internal/utils/response"
	"cornyk/gin-template/pkg/logger"
	"github.com/gin-gonic/gin"
)

// GetUsers 获取所有用户
func GetUsers(c *gin.Context) {

	logger.GetLogger(c, "debug").Info("This is an info log")

	users, err := daos.GetAllUsers(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	response.SucJson(c, response.Pagination{
		List:  users,
		Count: len(users),
	})
}
