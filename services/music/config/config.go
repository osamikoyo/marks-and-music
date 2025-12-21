package config

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	DefaultAddr        = "localhost:50052"
	DefaultMetricsAddr = "localhost:8080"
)

type Config struct {
	Addr        string `yaml:"addr" mapstructure:"addr"`
	MetricsAddr string `yaml:"metrics_addr" mapstructure:"metrics_addr"`

	SearchRequestTimeout time.Duration `yaml:"search_request_timeout" mapstructure:"search_request_timeout"`
	RepositoryTimeout    time.Duration `yaml:"repo_timeout" mapstructure:"repo_timeout"`

	Cache    CacheConfig    `yaml:"cache" mapstructure:"cache"`
	Postgres PostgresConfig `yaml:"postgres" mapstructure:"postgres"`
}

type CacheConfig struct {
	ExpTime                  time.Duration `yaml:"exp_time" mapstructure:"default_exp_time"`
	ExpiredItemsPurgeTimeout time.Duration `yaml:"exp_items_purge_timeout" mapstructure:"exp_items_purge_timeout"`
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
	ConnMaxLifetime int `yaml:"conn_max_lifetime_minutes" mapstructure:"conn_max_lifetime_minutes"`
}

func (c *PostgresConfig) GetDSN() (string, error) {
	if c.DSN != "" {
		return c.DSN, nil
	}

	if c.Host == "" {
		return "", fmt.Errorf("postgres host is required")
	}
	if c.User == "" {
		return "", fmt.Errorf("postgres user is required")
	}
	if c.DBName == "" {
		return "", fmt.Errorf("postgres dbname is required")
	}

	port := c.Port
	if port == 0 {
		port = 5432
	}

	var authPart string
	if c.Password != "" {
		authPart = c.User + ":" + url.QueryEscape(c.Password)
	} else {
		authPart = c.User
	}

	hostPort := c.Host
	if !strings.Contains(hostPort, ":") && port != 5432 {
		hostPort = hostPort + ":" + strconv.Itoa(port)
	} else if port != 5432 {
		hostPort = c.Host + ":" + strconv.Itoa(port)
	}

	base := fmt.Sprintf("postgres://%s@%s/%s", authPart, hostPort, c.DBName)

	values := url.Values{}

	if c.SSLMode != "" {
		values.Add("sslmode", c.SSLMode)
	} else {
		values.Add("sslmode", "disable")
	}

	if len(values) > 0 {
		base += "?" + values.Encode()
	}

	return base, nil
}

func NewConfig(path string, logger *logger.Logger) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.SetDefault("addr", DefaultAddr)
	v.SetDefault("metrics_addr", DefaultMetricsAddr)

	v.SetDefault("repo_timeout", 30*time.Second)
	v.SetDefault("search_request_timeout", 30*time.Second)

	v.SetDefault("cache.default_exp_time", 5*time.Minute)
	v.SetDefault("cahce.exp_items_purge_timeout", 10*time.Minute)

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

	v.BindEnv("repo_timeout", "APP_REPO_TIMEOUT")
	v.BindEnv("search_request_timeout", "APP_SEARCH_REQUEST_TIMEOUT")

	v.BindEnv("cache.default_exp_time", "APP_EXP_TIME")
	v.BindEnv("cache.exp_items_purge_timeout", "APP_EXP_ITEMS_PURGE_TIMEOUT")

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
