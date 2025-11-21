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

	Postgres PostgresConfig `yaml:"postgres" mapstructure:"postgres"`
}

type PostgresConfig struct {
	DSN string `yaml:"dsn" mapstructure:"dsn"`

	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
	DBName   string `yaml:"dbname" mapstructure:"dbname"`
	SSLMode  string `yaml:"sslmode" mapstructure:"sslmode"`

	MaxOpenConns    int `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	MaxIdleConns    int `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	ConnMaxLifetime int `yaml:"conn_max_lifetime_minutes" mapstructure:"conn_max_lifetime_minutes"` // в минутах
}

func NewConfig(path string, logger *logger.Logger) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.SetDefault("addr", DefaultAddr)
	v.SetDefault("metrics_addr", DefaultMetricsAddr)

	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", 5432)
	v.SetDefault("postgres.sslmode", "disable")
	v.SetDefault("postgres.max_open_conns", 25)
	v.SetDefault("postgres.max_idle_conns", 25)
	v.SetDefault("postgres.conn_max_lifetime_minutes", 5)

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

	v.BindEnv("postgres.dsn", "APP_POSTGRES_DSN")
	v.BindEnv("postgres.host", "APP_POSTGRES_HOST")
	v.BindEnv("postgres.port", "APP_POSTGRES_PORT")
	v.BindEnv("postgres.user", "APP_POSTGRES_USER")
	v.BindEnv("postgres.password", "APP_POSTGRES_PASSWORD")
	v.BindEnv("postgres.dbname", "APP_POSTGRES_DBNAME")
	v.BindEnv("postgres.sslmode", "APP_POSTGRES_SSLMODE")
	v.BindEnv("postgres.max_open_conns", "APP_POSTGRES_MAX_OPEN_CONNS")
	v.BindEnv("postgres.max_idle_conns", "APP_POSTGRES_MAX_IDLE_CONNS")
	v.BindEnv("postgres.conn_max_lifetime_minutes", "APP_POSTGRES_CONN_MAX_LIFETIME_MINUTES")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed unmarshal config: %w", err)
	}

	return &cfg, nil
}
