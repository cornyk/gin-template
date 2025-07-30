package main

import (
	"cornyk/gin-template/internal/exceptions"
	"cornyk/gin-template/internal/middlewares"
	"cornyk/gin-template/pkg/config"
	"cornyk/gin-template/pkg/database/mysql"
	"cornyk/gin-template/pkg/global"
	"cornyk/gin-template/pkg/logger"
	"cornyk/gin-template/pkg/queue/beanstalkd"
	"cornyk/gin-template/pkg/redis"
	"cornyk/gin-template/pkg/timezone"
	"cornyk/gin-template/routes"
	"fmt"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置文件并将配置文件内容保存到全局变量
	loadConfig := config.LoadConfig("config.yaml")
	global.GlobalConfig = loadConfig

	// 设置全局时区
	timezone.InitTimezone(loadConfig.App.Timezone)

	// 初始化日志
	logger.InitLogger()

	// 初始化MySQL
	mysql.InitDB(loadConfig)
	defer mysql.CloseAll()

	// 初始化Redis
	redis.InitRedis(loadConfig)
	defer redis.CloseAll()

	// 初始化Beanstalkd
	beanstalkd.InitBeanstalkd(loadConfig)
	defer beanstalkd.CloseAll()

	// 设置Gin模式并创建Gin路由
	if loadConfig.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// 设置全局中间件
	r.Use(middlewares.TraceIdMiddleware())
	r.Use(exceptions.ErrorHandler()) // 全局处理异常
	r.Use(middlewares.RequestLogMiddleware())

	// 设置路由
	routes.SetupRoutes(r)

	// 使用 endless 优雅启动服务器
	serverAddress := fmt.Sprintf("%s:%d", loadConfig.Server.Host, loadConfig.Server.Port)
	server := endless.NewServer(serverAddress, r)
	fmt.Println("\033[32m" + "Server started at: http://" + serverAddress + "\033[0m")
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
