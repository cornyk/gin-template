package main

import (
	"cornyk/gin-template/internal/middlewares"
	"cornyk/gin-template/pkg/config"
	"cornyk/gin-template/pkg/database"
	"cornyk/gin-template/pkg/global"
	"cornyk/gin-template/pkg/logger"
	"cornyk/gin-template/routes"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化日志
	logger.InitLogger()

	// 加载配置文件
	loadConfig, err := config.LoadConfig("config.yaml")
	if err != nil {
		panic("Failed to load config file")
	}

	// 将配置文件内容保存到全局变量
	global.GlobalConfig = loadConfig

	// 使用配置文件中的信息初始化数据库连接
	err = database.ConnectDB(loadConfig)
	if err != nil {
		panic("failed to connect to the database")
	}

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
	err = r.Run(serverAddress)
	if err != nil {
		fmt.Println("Server failed to start")
	}
}
