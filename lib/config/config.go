package config

import (
	"github.com/korableg/space-307-meetup/db"
	"github.com/korableg/space-307-meetup/lib/foo"
	"github.com/korableg/space-307-meetup/lib/rest"
)

type Config struct {
	Rest *rest.Config `yaml:"rest"`
	DB   *db.Config   `yaml:"db"`
	Foo  *foo.Config  `yaml:"foo"`
}

func NewConfig() *Config {
	return &Config{
		Rest: rest.NewConfig(),
		DB:   db.NewConfig("foodb://127.0.0.1:3324@db=main"),
		Foo: &foo.Config{
			Address: "",
		},
	}
}
