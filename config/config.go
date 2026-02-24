package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultPort  = "8080"
	defaultUIURL = "http://localhost:5173"
)

type Config struct {
	DatabaseURL  string
	Port         string
	BaseURL      string
	UIURL        string
	RollbarToken string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading: %v", err)
	}

	config := &Config{
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		Port:         os.Getenv("API_PORT"),
		BaseURL:      os.Getenv("BASE_URL"),
		UIURL:        os.Getenv("UI_URL"),
		RollbarToken: os.Getenv("ROLLBAR_TOKEN"),
	}

	if config.Port == "" {
		config.Port = defaultPort
	}

	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:" + config.Port
	}

	if config.UIURL == "" {
		config.UIURL = defaultUIURL
	}

	return config
}
