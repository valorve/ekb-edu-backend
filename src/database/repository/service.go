package repository

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.String()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (c Config) String() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		c.Host, c.User, c.Password, c.DB, c.Port)
}
