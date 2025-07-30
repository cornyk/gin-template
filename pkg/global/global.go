package global

import (
	"cornyk/gin-template/pkg/config"
	"cornyk/gin-template/pkg/queue/beanstalkd/connection"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// GlobalConfig 是全局配置结构体
var GlobalConfig *config.Config

// 定义全局数据库连接变量
var (
	DBConn         func(name ...string) *gorm.DB
	RedisConn      func(name ...string) *redis.Client
	BeanstalkdConn func(name ...string) *connection.TubeConn
)
