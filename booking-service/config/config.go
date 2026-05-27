package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPPort       string
	DatabaseURL    string
	JWTSecret      string
	JWTAccessTTL   time.Duration
	MigrationsPath string
}

func Load() *Config {
	return &Config{
		HTTPPort:       getEnv("HTTP_PORT", "8082"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/bookingdb?sslmode=disable"),
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
