package config

import (
	"errors"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

func New() (*Config, error) {
	config := Config{}

	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to read env"))
	}

	return &config, nil
}
