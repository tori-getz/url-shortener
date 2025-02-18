package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App  AppConfig
	Db   DbConfig
	Auth AuthConfig
}

type AppConfig struct {
	Addr string
}
type DbConfig struct {
	Dsn string
}

type AuthConfig struct {
	Secret string
}

func LoadConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		panic("Unable to load .env")
	}

	return &Config{
		App: AppConfig{
			Addr: os.Getenv("APP_ADDR"),
		},
		Db: DbConfig{
			Dsn: os.Getenv("DB_DSN"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("AUTH_SECRET"),
		},
	}
}
