package routes

import (
	"cornyk/gin-template/internal/commons/response_def"
	"cornyk/gin-template/internal/controllers"
	"cornyk/gin-template/internal/utils/response_util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置用户相关的路由
func SetupRoutes(router *gin.Engine) {

	userController := controllers.UserController{}

	router.GET("/users", userController.GetUsers)

	// health check url
	router.Any("/ping", func(c *gin.Context) {
		response_util.SucJson(c)
	})

	// 404 url
	router.NoRoute(func(c *gin.Context) {
		response_util.Json(c, response_def.CodeNoApi, response_def.MsgNoApi, nil, http.StatusNotFound)
	})
}
