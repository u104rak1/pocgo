package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Env struct {
	MYSQL_HOST     string `env:"MYSQL_HOST" envDefault:"localhost"`
	MYSQL_DBNAME   string `env:"MYSQL_DBNAME" envDefault:"POCGO_LOCAL_DB"`
	MYSQL_USER     string `env:"MYSQL_USER" envDefault:"local_user"`
	MYSQL_PASSWORD string `env:"MYSQL_PASSWORD" envDefault:"password"`
	MYSQL_PORT     string `env:"MYSQL_PORT" envDefault:"3306"`
	JWT_SECRET_KEY string `env:"JWT_SECRET_KEY" envDefault:"jwt_secret_key"`
}

func NewEnv() *Env {
	e, err := env.ParseAs[Env]()
	if err != nil {
		panic(fmt.Errorf("failed to parse env: %w", err))
	}
	return &e
}
