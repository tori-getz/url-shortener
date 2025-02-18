package main

import (
	"fmt"
	"net/http"
	config "url-shortener/configs"
	"url-shortener/internal/auth"
	"url-shortener/internal/link"
	"url-shortener/pkg/db"
)

func main() {
	cfg := config.LoadConfig()

	_ = db.NewDb(*cfg)

	router := http.NewServeMux()

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: cfg,
	})
	link.NewLinkHandler(router)

	server := http.Server{
		Addr:    cfg.App.Addr,
		Handler: router,
	}

	fmt.Println("App listen at ", cfg.App.Addr)
	server.ListenAndServe()
}
