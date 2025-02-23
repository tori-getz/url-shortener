package main

import (
	"os/user"
	config "url-shortener/configs"
	"url-shortener/internal/link"
	"url-shortener/pkg/db"
)

func main() {
	cfg := config.LoadConfig()
	database := db.NewDb(*cfg)
	database.AutoMigrate(&link.Link{})
	database.AutoMigrate(&user.User{})
}
