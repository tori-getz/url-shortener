package main

import (
	"fmt"
	"net/http"
	config "url-shortener/configs"
	"url-shortener/internal/auth"
	"url-shortener/internal/link"
	"url-shortener/pkg/db"
	"url-shortener/pkg/middleware"

	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()

	db := db.NewDb(*cfg)

	var logger *zap.Logger
	var err error

	if cfg.Logging.Mode == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err.Error())
	}

	router := http.NewServeMux()

	// Repositories
	linkRepository := link.NewLinkRepository(db)

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: cfg,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
	})

	middlewareChain := middleware.Compose(
		middleware.LoggingMiddleware(logger),
	)

	server := http.Server{
		Addr:    cfg.App.Addr,
		Handler: middlewareChain(router),
	}

	listenStr := fmt.Sprintf("App listen at %v", cfg.App.Addr)
	logger.Info(listenStr)
	server.ListenAndServe()
}
