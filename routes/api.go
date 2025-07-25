package routes

import (
	"cornyk/gin-template/internal/controllers"
	"cornyk/gin-template/internal/utils/response"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置用户相关的路由
func SetupRoutes(router *gin.Engine) {
	// health check url
	router.Any("/ping", func(c *gin.Context) {
		response.SucJson(c)
	})

	router.GET("/users", controllers.GetUsers)
}
