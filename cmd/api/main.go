package main

import (
	"cornyk/gin-template/internal/middlewares"
	"cornyk/gin-template/pkg/config"
	"cornyk/gin-template/pkg/database/mysql"
	"cornyk/gin-template/pkg/global"
	"cornyk/gin-template/pkg/logger"
	"cornyk/gin-template/pkg/redis"
	"cornyk/gin-template/routes"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化日志
	logger.InitLogger()

	// 加载配置文件并将配置文件内容保存到全局变量
	loadConfig := config.LoadConfig("config.yaml")
	global.GlobalConfig = loadConfig

	// 初始化MySQL
	mysql.InitDB(loadConfig)
	defer mysql.CloseAll()

	// 初始化Redis
	redis.InitRedis(loadConfig)
	defer redis.CloseAll()

	// 创建 Gin 路由
	r := gin.Default()

	// 设置全局中间件
	r.Use(middlewares.TraceIdMiddleware())
	r.Use(middlewares.RequestLogMiddleware())

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务器
	serverAddress := fmt.Sprintf("%s:%d", loadConfig.Server.Host, loadConfig.Server.Port)
	fmt.Println("\033[32m" + "Server started at: http://" + serverAddress + "\033[0m")
	err := r.Run(serverAddress)
	if err != nil {
		fmt.Println("Server failed to start")
	}
}
