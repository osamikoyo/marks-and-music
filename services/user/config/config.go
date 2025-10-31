package config

import "time"

const (
	DefaultAddr = "localhost:8081"
	DefaultJwtKey = "key"
	DefaultLogLevel = "debug"
	DefaultRTokenTTL = 72 * time.Hour
	DefaultATokenTTL = 1 * time.Minute
)

type Config struct {
	Addr     string `yaml:"addr"`
	JwtKey   string `yaml:"jwt_key"`
	LogLevel string `yaml:"log_level`
	RTokenTTL time.Duration `yaml:"refresh_token_ttl"`
	ATokenTTL time.Duration `yaml:"access_token_ttl"`
}

func NewConfig(path string) (*Config, error) {
	
}