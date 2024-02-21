package config

import (
	"ekb-edu/src/database/repository"
)

type Config struct {
	Web      Web
	Jwt      Jwt
	Postgres repository.Config
}

type Web struct {
	Port uint16 `env:"WEB_PORT" env-default:"8000"`
}

type Jwt struct {
	Secret string `env:"JWT_SECRET"`
}
