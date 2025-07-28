package redis

import (
	"context"
	"cornyk/gin-template/pkg/config"
	"cornyk/gin-template/pkg/global"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	clients     = make(map[string]*redis.Client)
	mu          sync.RWMutex
	ctx         = context.Background()
	defaultConn = "default" // 默认连接名称
)

// InitRedis 初始化Redis连接池
func InitRedis(config *config.Config) {
	mu.Lock()
	defer mu.Unlock()

	// 初始化所有配置的连接
	for name, cfg := range config.Redis {
		client := initConnection(cfg)
		if client != nil {
			clients[name] = client
		}
	}

	// 设置全局默认连接
	if _, exists := clients[defaultConn]; exists {
		global.RedisConn = func(name ...string) *redis.Client {
			return getClient(name...)
		}
	} else {
		panic("Default Redis connection not configured")
	}
}

func initConnection(cfg config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 添加前缀Hook
	if cfg.Prefix != "" {
		client.AddHook(NewPrefixHook(cfg.Prefix))
	}

	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Redis connect failed: %v", err))
	}

	return client
}

func getClient(names ...string) *redis.Client {
	name := defaultConn
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}

	mu.RLock()
	defer mu.RUnlock()

	if client, ok := clients[name]; ok {
		return client
	}
	panic(fmt.Sprintf("Redis connection [%s] uninitialized", name))
}

// CloseAll 关闭所有Redis连接
func CloseAll() {
	mu.Lock()
	defer mu.Unlock()

	for name, client := range clients {
		if err := client.Close(); err != nil {
			fmt.Printf("Close Redis connection [%s] failed: %v\n", name, err)
		}
		delete(clients, name)
	}
}
