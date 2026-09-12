package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() Config {
	// Load .env
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	if port == "" {
		port = "8000"
	}

	return Config{
		DatabaseURL: databaseURL,
		Port:        port,
	}
}