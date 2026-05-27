package config

import (
	"os"
	"time"
)

// Config holds all environment-driven settings for auth-service.
type Config struct {
	HTTPPort        string
	DatabaseURL     string
	JWTSecret       string
	JWTAccessTTL    time.Duration
	MigrationsPath  string
}

// Load reads configuration from environment variables with fallbacks.
func Load() *Config {
	return &Config{
		HTTPPort:       getEnv("HTTP_PORT", "8081"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/authdb?sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "super-secret-key-change-in-production"),
		JWTAccessTTL:   15 * time.Minute,
		MigrationsPath: getEnv("MIGRATIONS_PATH", "file://migrations"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
