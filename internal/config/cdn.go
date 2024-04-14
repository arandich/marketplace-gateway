package config

type CdnConfig struct {
	Host string `yaml:"cdn_host" env:"HOST" default:"0.0.0.0"`
	Port string `yaml:"cdn_port" env:"PORT" default:"2222"`
}
