// Package config 集中解析服务运行所需的全部环境变量配置。
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config 保存服务运行所需的全部配置项。
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Auth     AuthConfig
	Limiter  LimiterConfig
}

// ServerConfig HTTP 服务配置。
type ServerConfig struct {
	Port string
	Mode string
}

// DatabaseConfig MySQL 连接配置。
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// RedisConfig 缓存连接配置。
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// AuthConfig JWT 与 API Key 认证配置。
type AuthConfig struct {
	JWTSecret        string
	APIKeySecret     string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	APIKeyAdminName  string
}

// LimiterConfig 限流配置。
type LimiterConfig struct {
	PerMinute int
	Burst     int
}

// Load 从环境变量加载配置，缺失必填项时返回错误。
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("GIN_MODE", "release"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "supplychain_user"),
			Password: getEnv("DB_PASSWORD", "supplychain_pwd"),
			Name:     getEnv("DB_NAME", "supplychain_db"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Auth: AuthConfig{
			JWTSecret:       getEnv("JWT_SECRET", "change_me_to_a_long_random_string"),
			APIKeySecret:    getEnv("API_KEY_SECRET", "change_me_to_a_long_random_string"),
			AccessTokenTTL:  24 * time.Hour,
			RefreshTokenTTL: 7 * 24 * time.Hour,
			APIKeyAdminName: "api_key_service",
		},
		Limiter: LimiterConfig{
			PerMinute: getEnvInt("RATE_LIMIT_PER_MINUTE", 600),
			Burst:     getEnvInt("RATE_LIMIT_BURST", 1200),
		},
	}
	if cfg.Auth.JWTSecret == "" || cfg.Auth.APIKeySecret == "" {
		return nil, fmt.Errorf("load config: JWT_SECRET and API_KEY_SECRET must not be empty")
	}
	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

// DSN 返回 MySQL 驱动连接串。
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.Name)
}
