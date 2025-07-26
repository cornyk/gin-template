package config

import (
	"github.com/spf13/viper"
	"log"
)

// Config 存储配置信息
type Config struct {
	Server            ServerConfig   `yaml:"server"`
	Database          DatabaseConfig `yaml:"database"`
	SecondaryDatabase DatabaseConfig `yaml:"secondaryDatabase"`
	Redis             RedisConfig    `yaml:"redis"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	DBName    string `yaml:"dbname"`
	Charset   string `yaml:"charset"`
	ParseTime bool   `yaml:"parseTime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
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
