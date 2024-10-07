package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Env struct {
	JWT_SECRET_KEY string `env:"JWT_SECRET_KEY" envDefault:"jwt_secret_key"`
}

func NewEnv() *Env {
	e, err := env.ParseAs[Env]()
	if err != nil {
		panic(fmt.Errorf("failed to parse env: %w", err))
	}
	return &e
}
