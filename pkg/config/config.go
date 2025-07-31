package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config 存储配置信息
type Config struct {
	Server     ServerConfig                `yaml:"server"`
	App        AppConfig                   `yaml:"app"`
	Database   map[string]DatabaseConfig   `yaml:"database"`
	Redis      map[string]RedisConfig      `yaml:"redis"`
	Beanstalkd map[string]BeanstalkdConfig `yaml:"beanstalkd"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name     string `yaml:"name"`
	Env      string `yaml:"env"`
	Debug    bool   `yaml:"debug"`
	Timezone string `yaml:"timezone"`
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

// BeanstalkdConfig Beanstalkd 配置
type BeanstalkdConfig struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Tube    string `yaml:"tube"`
	Pri     uint32 `yaml:"pri"`
	Delay   int    `yaml:"delay"`   // 秒
	TTR     int    `yaml:"ttr"`     // 秒
	TimeOut int    `yaml:"timeout"` // 秒
}

// LoadConfig 读取配置文件并解析
func LoadConfig(configPath string) *Config {
	viper.SetConfigFile(configPath)

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic("Error reading config file, %s" + err.Error())
	}

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		panic("Unable to decode config file into struct, " + err.Error())
	}

	return &config
}
