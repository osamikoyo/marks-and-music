package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	DefaultAddr      = "localhost:8081"
	DefaultJwtKey    = "super-secret-jwt-key-change-in-production"
	DefaultLogLevel  = "debug"
	DefaultRTokenTTL = 72 * time.Hour
	DefaultATokenTTL = 15 * time.Minute // Рекомендуется 15 минут для access token
)

type Config struct {
	Addr      string        `yaml:"addr" mapstructure:"addr"`
	JwtKey    string        `yaml:"jwt_key" mapstructure:"jwt_key"`
	LogLevel  string        `yaml:"log_level" mapstructure:"log_level"`
	RTokenTTL time.Duration `yaml:"refresh_token_ttl" mapstructure:"refresh_token_ttl"`
	ATokenTTL time.Duration `yaml:"access_token_ttl" mapstructure:"access_token_ttl"`
}

func NewConfig(path string, logger *zap.Logger) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.SetDefault("addr", DefaultAddr)
	v.SetDefault("jwt_key", DefaultJwtKey)
	v.SetDefault("log_level", DefaultLogLevel)
	v.SetDefault("refresh_token_ttl", DefaultRTokenTTL)
	v.SetDefault("access_token_ttl", DefaultATokenTTL)

	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		logger.Warn("Config file not found, using defaults and environment variables", zap.String("path", path))
	} else {
		logger.Info("Config loaded from file", zap.String("path", path))
	}

	_ = v.BindEnv("addr", "APP_ADDR")
	_ = v.BindEnv("jwt_key", "APP_JWT_KEY")
	_ = v.BindEnv("log_level", "APP_LOG_LEVEL")
	_ = v.BindEnv("refresh_token_ttl", "APP_REFRESH_TOKEN_TTL")
	_ = v.BindEnv("access_token_ttl", "APP_ACCESS_TOKEN_TTL")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	logger.Info("Configuration loaded",
		zap.String("addr", cfg.Addr),
		zap.String("log_level", cfg.LogLevel),
		zap.Duration("access_token_ttl", cfg.ATokenTTL),
		zap.Duration("refresh_token_ttl", cfg.RTokenTTL),
		zap.Bool("jwt_key_set", cfg.JwtKey != DefaultJwtKey),
	)

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Addr == "" {
		return fmt.Errorf("addr is required")
	}

	if c.JwtKey == "" || c.JwtKey == DefaultJwtKey {
		return fmt.Errorf("jwt_key must be set and not default")
	}

	if c.ATokenTTL <= 0 {
		return fmt.Errorf("access_token_ttl must be positive")
	}

	if c.RTokenTTL <= 0 {
		return fmt.Errorf("refresh_token_ttl must be positive")
	}

	if c.RTokenTTL < time.Hour {
		return fmt.Errorf("refresh_token_ttl should be at least 1 hour")
	}

	if c.ATokenTTL > time.Hour {
		return fmt.Errorf("access_token_ttl should not exceed 1 hour")
	}

	return nil
}

