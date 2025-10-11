package config

import (
	"log"
	"os"
)

type Config struct {
	ServerAddress string
	DatabaseURL   string
	JWTSecret     string
	Environment   string
}

func Load() *Config {
	cfg := &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/project_budget_tracker?sslmode=disable"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		Environment:   getEnv("ENVIRONMENT", "development"),
	}

	log.Printf("Configuration loaded: Environment=%s, ServerAddress=%s", cfg.Environment, cfg.ServerAddress)
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
