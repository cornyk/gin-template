package logger

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// Logger 是我们全局的日志记录器
var Logger *logrus.Logger

// CustomFormatter 自定义日志格式
type CustomFormatter struct{}

// Format 实现 logrus.Formatter 接口
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 获取当前时间并格式化为 [YYYY-mm-dd HH:ii:ss] 格式
	logTime := time.Now().Format("2006-01-02 15:04:05")

	// 获取 traceId（如果存在）
	traceID, ok := entry.Data["trace-id"].(string)
	if !ok {
		traceID = "N/A" // 如果没有 trace-id，默认为 "N/A"
	}

	// 格式化日志内容为 [YYYY-mm-dd HH:ii:ss][traceId][logLevel] log Content
	logMessage := fmt.Sprintf("[%s][%s][%s]%s\n", logTime, traceID, entry.Level.String(), entry.Message)

	return []byte(logMessage), nil
}

// InitLogger 初始化日志配置
func InitLogger() {
	// 初始化日志记录器
	Logger = logrus.New()

	// 设置自定义日志格式
	Logger.SetFormatter(&CustomFormatter{})

	// 设置日志等级
	Logger.SetLevel(logrus.InfoLevel)
}

// GetLogger 获取全局日志实例，支持 channel 参数
func GetLogger(ctx context.Context, channel ...string) *logrus.Entry {
	// 如果没有传入 channel 参数，使用默认值 "app"
	if len(channel) == 0 {
		channel = append(channel, "app")
	}

	// 获取第一个 channel 参数
	logChannel := channel[0]

	// 获取 trace-id（如果存在），并将其传递给日志
	traceID := "N/A" // 默认值
	if traceIDValue, ok := ctx.Value("trace-id").(string); ok {
		traceID = traceIDValue
	}

	// 根据 channel 创建不同的日志输出，并将 trace-id 传递给 entry
	log := Logger.WithField("channel", logChannel).WithField("trace-id", traceID)

	// 设置不同的输出文件
	var logFile *os.File
	var err error

	// 获取当前日期，格式化为 YYYYmmdd
	currentDate := time.Now().Format("20060102")

	// 根据不同类型的日志文件来设置路径
	logPath := fmt.Sprintf("runtime/logs/%s-%s.log", logChannel, currentDate)

	// 如果日志文件不存在，创建它
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logger.Fatal("Error opening log file", err)
	}

	// 将日志输出到指定文件
	log.Logger.SetOutput(logFile)

	return log
}
