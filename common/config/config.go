package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Database struct {
	URI string `env:"DATABASE_URI"`
}

type Server struct {
	Addr         string        `env:"SERVER_ADDR"`
	ReadTimeout  time.Duration `env:"SERVER_READTIMEOUT" envDefault:"10s"`
	WriteTimeout time.Duration `env:"SERVER_WRITETIMEOUT" envDefault:"10s"`
}
type Config struct {
	DebugMode bool `env:"DEBUG"`

	Server Server
	DB     Database
}

func New() *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}
