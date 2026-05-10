package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port     int
	JWTToken string
}

func Load() (*Config, error) {
	port, err := envInt("PORT", 8080)
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:     port,
		JWTToken: envString("JWT_TOKEN", "secret"),
	}, nil
}

func envString(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}

	return value
}

func envInt(name string, defaultValue int) (int, error) {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be integer: %w", name, err)
	}

	if parsed <= 0 {
		return 0, fmt.Errorf("%s must be positive", name)
	}

	return parsed, nil
}
