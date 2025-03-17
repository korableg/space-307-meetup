package db

import (
	"time"
)

type Config struct {
	DSN          string        `yaml:"dsn"`
	Pool         int           `yaml:"pool"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

func NewConfig(dsn string) *Config {
	return &Config{
		Pool:         2,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 3 * time.Second,
		DSN:          dsn,
	}
}
