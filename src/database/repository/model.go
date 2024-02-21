package repository

type Config struct {
	Host     string `env:"POSTGRES_HOST"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Port     uint16 `env:"POSTGRES_PORT"`
	DB       string `env:"POSTGRES_DB"`
}
