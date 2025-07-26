package routes

import (
	"cornyk/gin-template/internal/controllers"
	"cornyk/gin-template/internal/utils/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SetupRoutes 设置用户相关的路由
func SetupRoutes(router *gin.Engine) {

	router.GET("/users", controllers.GetUsers)

	// health check url
	router.Any("/ping", func(c *gin.Context) {
		response.SucJson(c)
	})

	// 404 url
	router.NoRoute(func(c *gin.Context) {
		response.Json(c, response.CodeNoApi, response.MsgNoApi, nil, http.StatusNotFound)
	})
}
