package config

import (
	"time"

	"github.com/korableg/space-307-meetup/db"
	"github.com/korableg/space-307-meetup/lib/foo"
	"github.com/korableg/space-307-meetup/lib/rest"
)

type Config struct {
	Rest rest.Config `yaml:"rest"`
	DB   db.Config   `yaml:"db"`
	Foo  foo.Config  `yaml:"foo"`
}

func NewConfig() Config {
	return Config{
		Rest: rest.Config{
			Address: ":9090",
			Timeout: rest.ConfigTimeout{
				Handler: 5 * time.Second,
				Read:    time.Second,
			},
		},
		DB: db.Config{
			DSN: "foodb://127.0.0.1:3324@db=main",
		},
		Foo: foo.Config{
			Address: "",
		},
	}
}
