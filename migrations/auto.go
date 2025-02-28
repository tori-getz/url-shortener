package main

import (
	config "url-shortener/configs"
	"url-shortener/internal/link"
	"url-shortener/internal/stat"
	"url-shortener/internal/user"
	"url-shortener/pkg/db"
)

func main() {
	cfg := config.LoadConfig()
	database := db.NewDb(*cfg)
	database.AutoMigrate(&link.Link{}, &user.User{}, &stat.Stat{})
}
