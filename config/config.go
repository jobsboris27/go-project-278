package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
	BaseURL     string
}

func Load() *Config {
	config := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
		BaseURL:     os.Getenv("BASE_URL"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:" + config.Port
	}

	return config
}
