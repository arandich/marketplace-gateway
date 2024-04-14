package config

import "time"

type GoodsClientConfig struct {
	ConnString  string        `yaml:"conn_string" env:"CONN_STRING"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" default:"2m"`
}
