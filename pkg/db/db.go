package db

import (
	config "url-shortener/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func NewDb(cfg config.Config) *Db {
	db, err := gorm.Open(postgres.Open(cfg.Db.Dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	return &Db{
		DB: db,
	}
}
