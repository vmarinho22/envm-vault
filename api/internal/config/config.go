package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	Port        string
	DatabaseURL string
	JWTSecret   string
	JWTExpiry   time.Duration
	CORSOrigins []string
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

	origins, err := parseCORSOrigins(os.Getenv("CORS_ORIGINS"))
	if err != nil {
		return nil, err
	}
	cfg.CORSOrigins = origins

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

func parseCORSOrigins(raw string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		return []string{"http://localhost:3100"}, nil
	}

	var origins []string
	if err := json.Unmarshal([]byte(raw), &origins); err != nil {
		return nil, fmt.Errorf("parse CORS_ORIGINS: %w", err)
	}
	if len(origins) == 0 {
		return nil, fmt.Errorf("CORS_ORIGINS must contain at least one origin")
	}
	return origins, nil
}
