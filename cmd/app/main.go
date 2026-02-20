// Package main ...
package main

import (
	"flag"
	"log"
	"log/slog"
	"orgstructure/internal/config"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config.toml", "path to config file")
}

func main() {
	cfg := config.NewConfig()
	_, err := toml.DecodeFile(configPath, cfg)
	if err != nil {
		log.Fatal(err)
	}
	logger := config.NewLogger(cfg)
	log := logger.With(
		slog.String("BindAddr:", cfg.BindAddr),
		slog.String("DatabaseURL:", cfg.DatabaseURL),
		slog.String("LogLevel:", cfg.LogLevel),
		slog.String("TestDatabaseURL:", cfg.TestDatabaseURL),
	)
	log.Info("init")

}
