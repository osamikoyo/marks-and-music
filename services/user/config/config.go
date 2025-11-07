package config

import (
	"fmt"
	"time"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	DefaultAddr         = "localhost:8081"
	DefaultMetricsAddr = "localhost:8080"
	DefaultJwtKey       = "super-secret-jwt-key-change-in-production"
	DefaultRTokenTTL    = 72 * time.Hour
	DefaultATokenTTL    = 15 * time.Minute
	DefaultDatabasePath = "storage/users.db"
)

type Config struct {
	Addr         string        `yaml:"addr" mapstructure:"addr"`
	JwtKey       string        `yaml:"jwt_key" mapstructure:"jwt_key"`
	MetricsAddr  string        `yaml:"metrics_addr" mapstructure:"metrics_addr"`
	RTokenTTL    time.Duration `yaml:"refresh_token_ttl" mapstructure:"refresh_token_ttl"`
	ATokenTTL    time.Duration `yaml:"access_token_ttl" mapstructure:"access_token_ttl"`
	DatabasePath string        `yaml:"database_path" mapstructure:"database_path"`
}

func NewConfig(path string, logger *logger.Logger) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.SetDefault("addr", DefaultAddr)
	v.SetDefault("metrics_addr", DefaultMetricsAddr)
	v.SetDefault("jwt_key", DefaultJwtKey)
	v.SetDefault("refresh_token_ttl", DefaultRTokenTTL)
	v.SetDefault("access_token_ttl", DefaultATokenTTL)
	v.SetDefault("database_path", DefaultDatabasePath)

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
	_ = v.BindEnv("metrics_addr", "APP_METRICS_ADDR")
	_ = v.BindEnv("refresh_token_ttl", "APP_REFRESH_TOKEN_TTL")
	_ = v.BindEnv("access_token_ttl", "APP_ACCESS_TOKEN_TTL")
	_ = v.BindEnv("database_path", "APP_DATABASE_PATH")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	logger.Info("Configuration loaded",
		zap.String("addr", cfg.Addr),
		zap.String("metrics_addr", cfg.MetricsAddr),
		zap.Duration("access_token_ttl", cfg.ATokenTTL),
		zap.Duration("refresh_token_ttl", cfg.RTokenTTL),
		zap.Bool("jwt_key_set", cfg.JwtKey != DefaultJwtKey),
		zap.String("database_path", cfg.DatabasePath),
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

	if c.DatabasePath == "" {
		return fmt.Errorf("database_path should not be empty")
	}
 
	if c.Addr == c.MetricsAddr {
		return fmt.Errorf("metrics addr and grpc addr should be different")
	}

	return nil
}
