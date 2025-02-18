package main

import (
	"fmt"
	"net/http"
	config "url-shortener/configs"
	"url-shortener/internal/auth"
	"url-shortener/pkg/db"
)

func main() {
	cfg := config.LoadConfig()

	_ = db.NewDb(*cfg)

	router := http.NewServeMux()
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: cfg,
	})

	server := http.Server{
		Addr:    cfg.App.Addr,
		Handler: router,
	}

	fmt.Println("App listen at ", cfg.App.Addr)
	server.ListenAndServe()
}
