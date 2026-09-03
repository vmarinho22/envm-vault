package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	Port        string
	DatabaseURL string
	JWTSecret   string
	JWTExpiry   time.Duration
}

func Load() (*Config, error) {
	loadEnvFiles()

	jwtExpiry, err := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	if err != nil {
		return nil, fmt.Errorf("parse JWT_EXPIRY: %w", err)
	}

	cfg := &Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		JWTExpiry:   jwtExpiry,
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

func loadEnvFiles() {
	if err := godotenv.Load(); err == nil {
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		return
	}

	for {
		path := filepath.Join(dir, ".env")
		if _, statErr := os.Stat(path); statErr == nil {
			_ = godotenv.Load(path)
			return
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return
		}
		dir = parent
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
