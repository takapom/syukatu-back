package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL        string
	Port               string
	JWTSecret          string
	CORSAllowedOrigins []string
	Environment        string
}

func LoadConfig() *Config {
	// Load .env file if it exists (for local development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		DatabaseURL:  getEnv("DATABASE_URL", "sqlite://./example.db"),
		Port:         getEnv("PORT", "8080"),
		JWTSecret:    getEnv("JWT_SECRET", "your_dev_secret_key_which_is_long_enough"),
		Environment:  getEnv("ENVIRONMENT", "development"),
	}

	// Parse CORS allowed origins
	corsOrigins := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
	config.CORSAllowedOrigins = strings.Split(corsOrigins, ",")
	for i := range config.CORSAllowedOrigins {
		config.CORSAllowedOrigins[i] = strings.TrimSpace(config.CORSAllowedOrigins[i])
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}