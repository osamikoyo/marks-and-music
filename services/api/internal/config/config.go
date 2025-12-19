package config

type Config struct{
	Addr string `yaml:"addr" mapstructure:"addr"`
	MarkServiceAddr string `yaml:"mark_service_addr" mapstructure:"mark_service_addr"`
	MusicServiceAddr string `yaml:"music_service_addr" mapstructure:"music_service_addr"`
	UserServiceAddr string `yaml:"user_service_addr" mapstructure:"user_service_addr"`
}