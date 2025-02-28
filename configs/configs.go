package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App     AppConfig
	Logging LoggingConfig
	Db      DbConfig
	Auth    AuthConfig
}

type AppConfig struct {
	Addr string
	Logs string
}

type LoggingConfig struct {
	Mode string
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
		Logging: LoggingConfig{
			Mode: os.Getenv("LOGGING_MODE"),
		},
		Db: DbConfig{
			Dsn: os.Getenv("DB_DSN"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("AUTH_SECRET"),
		},
	}
}
