// Package config ...
package config

// Config ...
type Config struct {
	DatabaseURL     string `toml:"database_url"`
	TestDatabaseURL string `toml:"test_database_url"`
	BindAddr        string `toml:"bind_addr"`
	LogLevel        string `toml:"log_level"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "info",
	}
}
