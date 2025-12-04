package config

import (
	"fmt"
	"time"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Addr        string `yaml:"addr" mapstructure:"addr"`
	MetricsAddr string `yaml:"metrics_addr" mapstructure:"metrics_addr"`

	RepoTimeout time.Duration `yaml:"repo_timeout" mapstructure:"repo_timeout"`

	DBAddr string      `yaml:"db_addr" mapstructure:"db_addr"`

	Cache  CacheConfig `yaml:"cache" mapstrucure:"cache"`
}

type CacheConfig struct {
	ExpTime                  time.Duration `yaml:"exp_time" mapstructure:"default_exp_time"`
	ExpiredItemsPurgeTimeout time.Duration `yaml:"exp_items_purge_timeout" mapstructure:"exp_items_purge_timeout"`
}

func NewConfig(path string, logger *logger.Logger) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.SetDefault("addr", "localhost:50053")
	v.SetDefault("metrics_addr", "localhost:8083")

	v.SetDefault("repo_timeout", 30*time.Second)

	v.SetDefault("db_addr", "storage/marks.db")

	v.SetDefault("cache.default_exp_time", 5*time.Minute)
	v.SetDefault("cache.exp_items_purge_timeout", 10*time.Minute)

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
	v.BindEnv("metrics_addr", "APP_METRICS_ADDR")

	v.BindEnv("repo_timeout", "APP_REPO_TIMEOUT")

	v.BindEnv("db_addr", "APP_DB_ADDR")

	v.BindEnv("cache.default_exp_time", "APP_CACHE_DEFAULT_EXP_TIME")
	v.BindEnv("cache.exp_times_purge_timeout", "APP_CACHE_EXP_ITEMS_PURGE_TIMEOUT")

	var cfg Config
	if err := v.Unmarshal(&cfg);err != nil{
		return nil, fmt.Errorf("failed unmarshal config: %w", err)
	}

	return &cfg, nil
}
