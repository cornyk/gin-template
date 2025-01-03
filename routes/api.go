package routes

import (
	"cornyk/gin-template/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置用户相关的路由
func SetupRoutes(router *gin.Engine) {
	router.GET("/users", controllers.GetUsers)
}
