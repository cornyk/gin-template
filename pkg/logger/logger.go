package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	defaultTraceID = "N/A"
	defaultChannel = "app"
	logFilePerm    = 0666
	logDir         = "runtime/logs"
	dateFormat     = "20060102"
	timeFormat     = "2006-01-02 15:04:05"
	maxLogDays     = 7
	cleanInterval  = 24 * time.Hour
)

var (
	globalLogger  *logrus.Logger
	fileHandles   = make(map[string]*logFileHandle)
	fileHandleMux sync.RWMutex
	closeChan     = make(chan struct{})
)

type logFileHandle struct {
	file     *os.File
	lastUsed time.Time
	channel  string
}

type CustomFormatter struct {
	channel string
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	traceID := defaultTraceID
	if id, ok := entry.Data["trace-id"].(string); ok {
		traceID = id
	}

	return []byte(fmt.Sprintf("[%s][%s][%s]%s\n",
		time.Now().Format(timeFormat),
		traceID,
		entry.Level.String(),
		entry.Message)), nil
}

func InitLogger() {
	globalLogger = logrus.New()
	globalLogger.SetLevel(logrus.InfoLevel)

	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
	}

	go startCleanupTask()
}

func GetLogger(ctx context.Context, channel ...string) *logrus.Entry {
	logChannel := defaultChannel
	if len(channel) > 0 {
		logChannel = channel[0]
	}

	traceID := defaultTraceID
	if traceIDValue, ok := ctx.Value("trace-id").(string); ok {
		traceID = traceIDValue
	}

	logger := logrus.New()
	logger.SetFormatter(&CustomFormatter{channel: logChannel})
	logger.SetLevel(globalLogger.Level)

	handle := getFileHandle(logChannel)
	logger.SetOutput(handle.file)

	return logger.WithFields(logrus.Fields{
		"channel":  logChannel,
		"trace-id": traceID,
	})
}

func getFileHandle(channel string) *logFileHandle {
	fileHandleMux.Lock()
	defer fileHandleMux.Unlock()

	currentDate := time.Now().Format(dateFormat)
	fileKey := channel + "-" + currentDate

	if handle, exists := fileHandles[fileKey]; exists {
		handle.lastUsed = time.Now()
		return handle
	}

	logPath := filepath.Join(logDir, fileKey+".log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, logFilePerm)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return &logFileHandle{file: os.Stderr, channel: channel}
	}

	handle := &logFileHandle{
		file:     file,
		lastUsed: time.Now(),
		channel:  channel,
	}
	fileHandles[fileKey] = handle

	return handle
}

func Cleanup() {
	close(closeChan)
	fileHandleMux.Lock()
	defer fileHandleMux.Unlock()

	for key, handle := range fileHandles {
		if err := handle.file.Close(); err != nil {
			fmt.Printf("Failed to close log file %s: %v\n", key, err)
		}
		delete(fileHandles, key)
	}
}

func startCleanupTask() {
	ticker := time.NewTicker(cleanInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cleanOldLogs()
			cleanInactiveHandles()
		case <-closeChan:
			return
		}
	}
}

func cleanOldLogs() {
	files, err := os.ReadDir(logDir)
	if err != nil {
		fmt.Printf("Failed to read log directory: %v\n", err)
		return
	}

	cutoffTime := time.Now().AddDate(0, 0, -maxLogDays)
	var filesToDelete []string

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), ".log") {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoffTime) {
			filesToDelete = append(filesToDelete, file.Name())
		}
	}

	for _, filename := range filesToDelete {
		filePath := filepath.Join(logDir, filename)
		if err := os.Remove(filePath); err != nil {
			fmt.Printf("Failed to delete old log file %s: %v\n", filePath, err)
		}
	}
}

func cleanInactiveHandles() {
	fileHandleMux.Lock()
	defer fileHandleMux.Unlock()

	cutoffTime := time.Now().Add(-1 * time.Hour)
	for key, handle := range fileHandles {
		if handle.lastUsed.Before(cutoffTime) {
			_ = handle.file.Close()
			delete(fileHandles, key)
		}
	}
}
