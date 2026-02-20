package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
	BaseURL     string
	UIURL       string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading: %v", err)
	}

	config := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
		BaseURL:     os.Getenv("BASE_URL"),
		UIURL:       os.Getenv("UI_URL"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:" + config.Port
	}

	if config.UIURL == "" {
		config.UIURL = "http://localhost:5173"
	}

	return config
}
