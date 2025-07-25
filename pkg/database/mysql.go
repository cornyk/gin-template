package database

import (
	"cornyk/gin-template/pkg/config"
	"cornyk/gin-template/pkg/global"
	"cornyk/gin-template/pkg/logger"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// ConnectDB 连接到多个数据库，并保存到全局变量中
func ConnectDB(config *config.Config) error {
	// 主数据库连接
	mainDB, err := openDBConnection(config.Database)
	if err != nil {
		return err
	}
	global.MainDB = mainDB

	// 副本数据库连接
	secondaryDB, err := openDBConnection(config.SecondaryDatabase)
	if err != nil {
		return err
	}
	global.SecondaryDB = secondaryDB

	return nil
}

// 设置回调钩子，将 trace-id 从 Gin 的上下文传递到 Gorm 的上下文
func setTraceIDCallback() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		// 从 db.Statement.Context 中获取 trace-id
		traceID, _ := db.Statement.Context.Value("trace-id").(string)

		// 如果没有 trace-id，则使用默认值
		if traceID == "" {
			traceID = "N/A"
		}

		// 将 trace-id 设置到 Gorm 的 Statement 中
		db.Statement.Set("trace-id", traceID)

		// 确保可以正确记录日志
		sqlLogger := logger.GetLogger(db.Statement.Context)
		if sqlLogger != nil {
			sqlLogger.Info("Set trace-id in gorm statement", "trace-id", traceID)
		}
	}
}

// openDBConnection 根据配置建立数据库连接
func openDBConnection(dbConfig config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
		dbConfig.Charset,
		dbConfig.ParseTime,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.NewGormLogger(), // 使用自定义的 SQL Logger
	})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return nil, err
	}

	// 注册钩子，在每次数据库操作时传递 trace-id
	db.Callback().Create().Before("gorm:create").Register("set_trace_id", setTraceIDCallback())

	return db, nil
}
