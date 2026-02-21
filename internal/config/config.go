package config

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

// Config ...
type Config struct {
	DatabaseURL string `env:"DATABASE_URL,required"`
	BindAddr    string `env:"BIND_ADDR,default=:8080"`
	LogLevel    string `env:"LOG_LEVEL,default=info"`
}

// ParseConfig загружает .env файл (если есть) и парсит переменные окружения.
func ParseConfig() (Config, error) {
	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load env variables from file: %w", err)
	}

	var c Config
	if err := envconfig.Process(context.Background(), &c); err != nil {
		return Config{}, fmt.Errorf("parse env variables to config: %w", err)
	}

	return c, nil
}
