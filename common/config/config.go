package config

import (
	"encoding/base64"
	"errors"
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
type OAuth struct {
	ClientID string    `env:"CLIENT_ID,required"`
	Secret   SecretKey `env:"SECRET,required"`
}

type Server struct {
	Addr         string        `env:"SERVER_ADDR"`
	ReadTimeout  time.Duration `env:"SERVER_READTIMEOUT" envDefault:"10s"`
	WriteTimeout time.Duration `env:"SERVER_WRITETIMEOUT" envDefault:"10s"`
}
type Config struct {
	DebugMode bool `env:"DEBUG"`

	OAuth  OAuth
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

// SecretKey represents a 32-byte symmetric key decoded from base64.
type SecretKey []byte

// UnmarshalText lets env parse a base64-encoded string into a fixed-size key.
func (k *SecretKey) UnmarshalText(text []byte) error {
	decoded, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return err
	}
	if len(decoded) != 32 {
		return errors.New("secret must decode to 32 bytes")
	}
	*k = decoded
	return nil
}
