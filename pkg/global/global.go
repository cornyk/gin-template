package global

import (
	"cornyk/gin-template/pkg/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// GlobalConfig 是全局配置结构体
var GlobalConfig *config.Config

// 定义全局数据库连接变量

var MainDB *gorm.DB
var SecondaryDB *gorm.DB

var RedisConn func(names ...string) *redis.Client
