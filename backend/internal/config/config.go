package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config holds all application configurations
type Config struct {
	Environment string         `mapstructure:"environment"`
	LogLevel    string         `mapstructure:"log_level"`
	Server      ServerConfig   `mapstructure:"server"`
	Database    DatabaseConfig `mapstructure:"database"`
	Redis       RedisConfig    `mapstructure:"redis"`
	MinIO       MinIOConfig    `mapstructure:"minio"`
	NATS        NATSConfig     `mapstructure:"nats"`
	JWT         JWTConfig      `mapstructure:"jwt"`
	Wechat      WechatConfig   `mapstructure:"wechat"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         int `mapstructure:"port"`
	ReadTimeout  int `mapstructure:"read_timeout"`
	WriteTimeout int `mapstructure:"write_timeout"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// MinIOConfig holds MinIO S3 configuration
type MinIOConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	UseSSL    bool   `mapstructure:"use_ssl"`
}

// NATSConfig holds NATS JetStream configuration
type NATSConfig struct {
	URL string `mapstructure:"url"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// WechatConfig holds WeChat Pay configuration
type WechatConfig struct {
	MchID    string `mapstructure:"mch_id"`
	APIV3Key string `mapstructure:"api_v3_key"`
	AppID    string `mapstructure:"app_id"`
	CertPath string `mapstructure:"cert_path"`
	KeyPath  string `mapstructure:"key_path"`
}

// Load reads configuration from environment variables and config files
func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/backend/")

	// Set defaults
	viper.SetDefault("environment", "development")
	viper.SetDefault("log_level", "info")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("jwt.expire_hours", 24)
	viper.SetDefault("minio.bucket", "cruisebooking")
	viper.SetDefault("minio.use_ssl", false)

	// Enable environment variable override
	viper.AutomaticEnv()

	// Read config file (optional - env vars take precedence)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: config file not found: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}

	return &cfg
}

// GetDSN returns PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}
