package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

// Config 存储配置信息
type Config struct {
	Server   ServerConfig              `yaml:"server"`
	Database map[string]DatabaseConfig `yaml:"database"`
	Redis    map[string]RedisConfig    `yaml:"redis"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	DBName          string        `yaml:"dbname"`
	Charset         string        `yaml:"charset"`
	ParseTime       bool          `yaml:"parse_time"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`    // 默认10
	MaxOpenConns    int           `yaml:"max_open_conns"`    // 默认100
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"` // 默认1h
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"` // 默认100
	Prefix   string `yaml:"prefix"`
}

// LoadConfig 读取配置文件并解析
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return nil, err
	}

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
		return nil, err
	}

	return &config, nil
}
