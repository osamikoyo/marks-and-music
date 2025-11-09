package config

import (
	"fmt"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	DefaultAddr         = "localhost:8082"
	DefaultMetricsAddr  = "localhost:8080"
	DefaultDatabasePath = "storage/music.db"
)

type Config struct {
	Addr        string `yaml:"addr" mapstructure:"addr"`
	MetricsAddr string `yaml:"metrics_addr" mapstructure:"metrics_addr"`
	DBPath      string `yaml:"database_path" mapstructure:"database_path"`
}

func NewConfig(path string, logger *logger.Logger) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.SetDefault("addr", DefaultAddr)
	v.SetDefault("metrics_addr", DefaultMetricsAddr)
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

	v.BindEnv("addr", "APP_ADDR")
	v.BindEnv("metrics_addr", "APP_METRICS_ADDR")
	v.BindEnv("database_path", "APP_DATABASE_PATH")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed unmarshal config: %w", err)
	}

	return &cfg, nil
}

