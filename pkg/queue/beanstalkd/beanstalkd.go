package beanstalkd

import (
	"cornyk/gin-template/pkg/config"
	"cornyk/gin-template/pkg/global"
	"cornyk/gin-template/pkg/queue/beanstalkd/connection"
	"fmt"
)

// InitBeanstalkd 初始化所有 beanstalkd 连接
func InitBeanstalkd(config *config.Config) {
	// 初始化所有配置的 beanstalkd 连接
	for name, cfg := range config.Beanstalkd {
		tubeConn, err := connection.NewTubeConnection(cfg.Host, cfg.Port, cfg.Tube, cfg.Pri, cfg.Delay, cfg.TTR, cfg.TimeOut)
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to beanstalkd %s: %v", name, err))
		}
		connection.SetConnection(name, tubeConn)
	}

	// 设置全局连接函数
	global.BeanstalkdConn = func(names ...string) *connection.TubeConn {
		name := "default"
		if len(names) > 0 && names[0] != "" {
			name = names[0]
		}

		conn, ok := connection.GetConnection(name)
		if !ok {
			panic(fmt.Sprintf("Beanstalkd connection %s not found", name))
		}
		return conn
	}
}

// CloseAll 关闭所有 beanstalkd 连接
func CloseAll() {
	connection.CloseAll()
}
