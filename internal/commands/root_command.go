package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "cli",
	Short: "命令行工具",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 延迟获取Flag（此时cmd已初始化完成）
		configPath := cmd.Flags().Lookup("config").Value.String()
		initConfig(configPath)
		initDB()
	},
}

func init() {
	// 添加全局Flag（不直接绑定到initConfig）
	RootCmd.PersistentFlags().String("config", "./config.yaml", "配置文件路径")
	//RootCmd.AddCommand(dbCmd) // 确保子命令已定义
}

// 修改为接收参数
func initConfig(configPath string) {
	fmt.Printf("加载配置文件: %s\n", configPath)
	// 实际配置加载...
}

func initDB() {
	fmt.Println("初始化数据库连接")
}
