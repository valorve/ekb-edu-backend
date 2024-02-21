package storage

import "ekb-edu/src/database/repository"

func Connect(cfg *repository.Config) {
	var err error
	DB, err = repository.Connect(cfg)
	if err != nil {
		panic(err)
	}
}
