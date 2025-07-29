package middlewares

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TraceIdMiddleware 为每个请求生成 trace-id
func TraceIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成新的 trace-id
		traceID := strings.ReplaceAll(uuid.New().String(), "-", "")

		// 将 trace-id 添加到 Gin 上下文
		c.Set("trace-id", traceID)

		// 在 gin 上下文中传递 trace-id
		ctx := context.WithValue(c.Request.Context(), "trace-id", traceID)

		// 将修改后的 ctx 传递到 c.Request 上，以便后续处理中使用
		c.Request = c.Request.WithContext(ctx)

		// 执行请求
		c.Next()
	}
}
