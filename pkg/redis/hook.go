package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

// PrefixHook 自动为 Redis Key 添加前缀（prefix 为空时不处理）
type PrefixHook struct {
	prefix string
}

// NewPrefixHook 创建前缀 Hook，如果 prefix 为空则返回 nil（不生效）
func NewPrefixHook(prefix string) *PrefixHook {
	if prefix == "" {
		return nil // 如果 prefix 为空，直接返回 nil，不启用 Hook
	}
	return &PrefixHook{prefix: prefix}
}

// DialHook 连接钩子（无需修改）
func (h *PrefixHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

// ProcessHook 处理单个 Redis 命令
func (h *PrefixHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		h.processCommand(cmd) // 处理 Key
		return next(ctx, cmd)
	}
}

// ProcessPipelineHook 处理管道命令
func (h *PrefixHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		for _, cmd := range cmds {
			h.processCommand(cmd) // 处理每个命令的 Key
		}
		return next(ctx, cmds)
	}
}

// processCommand 实际处理 Key 的逻辑
func (h *PrefixHook) processCommand(cmd redis.Cmder) {
	// 只处理带 Key 的命令
	switch c := cmd.(type) {
	case *redis.StringCmd, *redis.StatusCmd, *redis.IntCmd, *redis.BoolCmd,
		*redis.SliceCmd, *redis.StringSliceCmd, *redis.FloatCmd, *redis.ScanCmd, *redis.ZSliceCmd:

		if len(c.Args()) > 1 {
			keyPos := 1 // 大多数命令的 Key 在 args[1]（如 GET/SET）
			switch cmd.(type) {
			case *redis.StringSliceCmd, *redis.ScanCmd:
				keyPos = 0 // 对于 SMEMBERS、SCAN 等命令，Key 在 args[0]
			}

			// 如果参数是 string 类型，并且非空，则添加前缀
			if key, ok := c.Args()[keyPos].(string); ok && key != "" {
				c.Args()[keyPos] = h.prefix + key // 直接拼接，不加 ":"
			}
		}
	}
}
