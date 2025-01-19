package rest

import "time"

type (
	ConfigTimeout struct {
		Handler time.Duration `yaml:"timeout"`
		Read    time.Duration `yaml:"read"`
	}

	Config struct {
		Address string        `yaml:"address"`
		Timeout ConfigTimeout `yaml:"timeout"`
	}
)
