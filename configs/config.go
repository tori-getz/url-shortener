package config

type Config struct {
	Db DbConfig
}

type DbConfig struct {
	Dsn string
}
