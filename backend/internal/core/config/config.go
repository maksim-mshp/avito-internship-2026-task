package config

import (
	"fmt"

	"go-simpler.org/env"
)

type Config struct {
	Port     int `env:"PORT" default:"8080"`
	Database Database
	JWTToken string `env:"JWT_TOKEN" default:"secret"`
}

type Database struct {
	Host     string `env:"POSTGRES_HOST" default:"localhost"`
	Port     int    `env:"POSTGRES_PORT" default:"5432"`
	User     string `env:"POSTGRES_USER" default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" default:"postgres"`
	Database string `env:"POSTGRES_DATABASE" default:"ai_assistants_catalog"`
}

func Load() (*Config, error) {
	cfg := Config{}
	if err := env.Load(&cfg, nil); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &cfg, nil
}
