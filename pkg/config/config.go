package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerHost string `envconfig:"SERVER_HOST"`
	ServerPort int    `envconfig:"SERVER_PORT"`
}

func ParseConfig() (*Config, error) {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		return nil, fmt.Errorf("could not process config: %w", err)
	}
	return &c, nil
}
