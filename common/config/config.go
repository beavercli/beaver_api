package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Database struct {
	URI               string        `env:"DATABASE_URI,required"`
	MaxConns          int32         `env:"DATABASE_MAX_CONNS" envDefault:"25"`
	MinConns          int32         `env:"DATABASE_MIN_CONNS" envDefault:"5"`
	MaxConnLifetime   time.Duration `env:"DATABASE_MAX_CONN_LIFETIME" envDefault:"1h"`
	MaxConnIdleTime   time.Duration `env:"DATABASE_MAX_CONN_IDLE_TIME" envDefault:"30m"`
	HealthCheckPeriod time.Duration `env:"DATABASE_HEALTH_CHECK_PERIOD" envDefault:"1m"`
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
