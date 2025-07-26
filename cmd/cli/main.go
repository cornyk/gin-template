package main

import (
	"cornyk/gin-template/internal/commands"
	"cornyk/gin-template/pkg/config"
	"cornyk/gin-template/pkg/database/mysql"
	"cornyk/gin-template/pkg/global"
	"cornyk/gin-template/pkg/logger"
	"github.com/spf13/cobra"
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

	// 根命令
	rootCmd := &cobra.Command{
		Use:   "cli",
		Short: "项目命令行工具",
	}

	// 注册子命令
	rootCmd.AddCommand(commands.MakeModelCmd())

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
