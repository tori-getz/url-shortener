package main

import (
	"fmt"
	"net/http"
	config "url-shortener/configs"
	"url-shortener/internal/auth"
	"url-shortener/internal/link"
	"url-shortener/internal/stat"
	"url-shortener/internal/user"
	"url-shortener/pkg/db"
	"url-shortener/pkg/event"
	"url-shortener/pkg/jwt"
	"url-shortener/pkg/middleware"

	_ "url-shortener/docs"

	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title URL Shortener API
// @version 1.0

// @contact.name Daniil Benger (tori-getz)
// @contact.url http://t.me/torigetz/
// @contact.email torigetz@yandex.ru

// @host localhost:3000

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

	// Event Bus
	eventBus := event.NewEventBus(logger)

	// Repositories
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	statRepository := stat.NewStatRepository(db)

	// services
	authService := auth.NewAuthService(*userRepository)
	statService := stat.NewStatService(stat.StatServiceDeps{
		StatRepository: statRepository,
		EventBus:       eventBus,
	})

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		AuthService: authService,
		Jwt:         jwt.NewJwt(cfg.Auth.Secret),
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		Config:         cfg,
		LinkRepository: linkRepository,
		EventBus:       eventBus,
	})
	stat.NewStatHandler(router, stat.StatHandlerDeps{
		Config:         cfg,
		StatRepository: statRepository,
	})

	middlewareChain := middleware.Compose(
		middleware.LoggingMiddleware(logger),
		middleware.CORS,
	)

	router.Handle("/docs/", httpSwagger.WrapHandler)

	server := http.Server{
		Addr:    cfg.App.Addr,
		Handler: middlewareChain(router),
	}

	go statService.AddClick()

	listenStr := fmt.Sprintf("App listen at %v", cfg.App.Addr)
	logger.Info(listenStr)
	server.ListenAndServe()
}
