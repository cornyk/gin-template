package mysql

import (
	"cornyk/gin-template/pkg/config"
	"cornyk/gin-template/pkg/global"
	"cornyk/gin-template/pkg/logger"
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbs       = make(map[string]*gorm.DB)
	mu        sync.RWMutex
	defaultDB = "default" // 默认连接名称
)

// InitDB 初始化MySQL连接池
func InitDB() {
	mu.Lock()
	defer mu.Unlock()

	// 初始化所有配置的连接
	for name, cfg := range global.GlobalConfig.Database {
		db := initConnection(cfg)
		dbs[name] = db
	}

	// 设置全局访问函数
	global.DBConn = func(name ...string) *gorm.DB {
		return getDB(name...)
	}
}

func initConnection(cfg config.DatabaseConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
		cfg.ParseTime,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.NewGormLogger(),
	})
	if err != nil {
		panic(fmt.Sprintf("MySQL connect failed: %v", err))
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to open MySQL connection: %v", err))
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// 注册回调
	registerCallbacks(db)

	return db
}

func getDB(names ...string) *gorm.DB {
	name := defaultDB
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}

	mu.RLock()
	defer mu.RUnlock()

	if db, ok := dbs[name]; ok {
		return db
	}
	panic(fmt.Sprintf("MySQL connection [%s] uninitialized", name))
}

// CloseAll 关闭所有数据库连接
func CloseAll() {
	mu.Lock()
	defer mu.Unlock()

	for name, db := range dbs {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
		delete(dbs, name)
	}
}

func registerCallbacks(db *gorm.DB) {
	_ = db.Callback().Create().Before("gorm:create").Register("set_trace_id", setTraceIDCallback())
	// 可以添加其他回调...
}

func setTraceIDCallback() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		traceID, _ := db.Statement.Context.Value("trace-id").(string)
		if traceID == "" {
			traceID = "N/A"
		}
		db.Statement.Set("trace-id", traceID)
	}
}
