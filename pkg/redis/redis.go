package redis

import (
	"context"
	"cornyk/gin-template/pkg/config"
	"cornyk/gin-template/pkg/global"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func ConnectRedis(config *config.Config) {
	mainRedisConn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	// 如果配置了 prefix，才启用 Hook
	if hook := NewPrefixHook(config.Redis.Prefix); hook != nil {
		mainRedisConn.AddHook(hook)
	}

	// 测试连接
	if err := mainRedisConn.Ping(ctx).Err(); err != nil {
		panic("MainRedis connect error: " + err.Error())
	}

	global.MainRedis = mainRedisConn
}
