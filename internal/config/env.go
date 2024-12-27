package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Env struct {
	APP_PORT          string `env:"APP_PORT" envDefault:"8080"`
	USE_INMEMORY      bool   `env:"USE_INMEMORY" envDefault:"false"`
	POSTGRES_HOST     string `env:"POSTGRES_HOST" envDefault:"postgres"`
	POSTGRES_DBNAME   string `env:"POSTGRES_DBNAME" envDefault:"POCGO_LOCAL_DB"`
	POSTGRES_USER     string `env:"POSTGRES_USER" envDefault:"local_user"`
	POSTGRES_PASSWORD string `env:"POSTGRES_PASSWORD" envDefault:"password"`
	POSTGRES_PORT     string `env:"POSTGRES_PORT" envDefault:"5432"`
	POSTGRES_SSLMODE  string `env:"POSTGRES_SSLMODE" envDefault:"disable"`
	JWT_SECRET_KEY    string `env:"JWT_SECRET_KEY" envDefault:"jwt_secret_key"`
}

func NewEnv() *Env {
	e, err := env.ParseAs[Env]()
	if err != nil {
		panic(fmt.Errorf("failed to parse env: %w", err))
	}
	return &e
}
