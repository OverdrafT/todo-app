package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Port       string `env:"SERVER_PORT" default:"8000"`
	LogLevel   string `env:"LOG_LEVEL" default:"INFO"`
	Repository repo   `env:"REPOSITORY" required:"true"`
	Username   string `env:"POSTGRES_USER"`
	Password   string `env:"POSTGRES_PASSWORD"`
	DB         string `env:"POSTGRES_DB"`
}

type repo string

const (
	MemoryRepo   repo = "memory"
	PostgresRepo repo = "postgres"
)

func Init() (cfg Config) {
	var err error
	if err = env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}
