package config

import (
	"fmt"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Addr             string `yaml:"addr" mapstructure:"addr"`
	MarkServiceAddr  string `yaml:"mark_service_addr" mapstructure:"mark_service_addr"`
	MusicServiceAddr string `yaml:"music_service_addr" mapstructure:"music_service_addr"`
	UserServiceAddr  string `yaml:"user_service_addr" mapstructure:"user_service_addr"`
}

func NewConfig(path string, logger *logger.Logger) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.SetDefault("addr", "localhost:8080")
	v.SetDefault("mark_service_addr", "localhost:50053")
	v.SetDefault("music_service_addr", "localhost:50052")
	v.SetDefault("user_service_addr", "localhost:50051")

	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		logger.Warn("Config file not found, using defaults and environment variables", zap.String("path", path))
	} else {
		logger.Info("Config loaded from file",
			zap.String("path", path))
	}

	v.BindEnv("addr", "APP_ADDR")
	v.BindEnv("user_service_addr", "APP_USER_SERVICE_ADDR")
	v.BindEnv("music_service_addr", "APP_MUSIC_SERVICE_ADDR")
	v.BindEnv("mark_service_addr", "APP_MARK_SERVICE_ADDR")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed unmarshal config: %w", err)
	}

	return &cfg, nil
}
