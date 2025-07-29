package logger

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
)

// NewGormLogger 创建一个定制的 Gorm Logger
func NewGormLogger() logger.Interface {
	return &gormLogger{}
}

type gormLogger struct{}

// LogMode 设置日志级别
func (g *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

// getTraceID 获取 trace-id，若未找到则返回默认值
func getTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value("trace-id").(string); ok && traceID != "" {
		return traceID
	}
	return "N/A" // 默认值
}

// formatLogMessage 格式化日志信息
func formatLogMessage(msg string, data ...interface{}) string {
	return fmt.Sprintf(msg, data...)
}

// logInfo 统一记录日志信息的函数
func (g *gormLogger) logInfo(ctx context.Context, logLevel, msg string, data ...interface{}) {
	logMessage := formatLogMessage(msg, data...)
	finalMessage := fmt.Sprintf("%s", logMessage)

	// 记录日志
	switch logLevel {
	case "info":
		GetLogger(ctx, "sql").Info(finalMessage)
	case "warn":
		GetLogger(ctx, "sql").Warn(finalMessage)
	case "error":
		GetLogger(ctx, "sql_err").Error(finalMessage)
	}
}

// Info 记录 Info 级别的日志
func (g *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	g.logInfo(ctx, "info", msg, data...)
}

// Warn 记录 Warn 级别的日志
func (g *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	g.logInfo(ctx, "warn", msg, data...)
}

// Error 记录 Error 级别的日志
func (g *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	g.logInfo(ctx, "error", msg, data...)
}

// Trace 记录 SQL 查询日志
func (g *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	// 格式化 SQL 日志
	sqlLog := fmt.Sprintf("[%v][rows:%v]%s", elapsed, rows, sql)

	// 如果有错误，记录到sql_err通道
	if err != nil {
		GetLogger(ctx, "sql_err").
			WithError(err).
			Error(fmt.Sprintf("%s - Error: %v", sqlLog, err))
	} else {
		GetLogger(ctx, "sql").Info(sqlLog)
	}
}
