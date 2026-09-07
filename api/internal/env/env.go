package env

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Envs is the central process environment for the Go API.
// Usage: env.Load() then env.Envs.Database.URL (same idea as `import { envs } from 'env'`).
//
// Required: DATABASE_URL, JWT_SECRET, ENCRYPTION_KEY.
// Optional (defaults): APP_ENV, PORT, JWT_EXPIRY, CORS_ORIGINS, RATE_LIMIT_*, ENCRYPTION_KEY_VERSION.
// POSTGRES_* are Docker Compose / .docker/db/.env only — not loaded here.
var Envs Vars

// Vars holds validated environment values after Load(), grouped by context.
type Vars struct {
	App       AppEnv
	Database  DatabaseEnv
	JWT       JWTEnv
	Crypto    CryptoEnv
	CORS      CORSEnv
	RateLimit RateLimitEnv
}

type AppEnv struct {
	Env  string // APP_ENV
	Port string
}

type DatabaseEnv struct {
	URL string // DATABASE_URL
}

type JWTEnv struct {
	Secret string
	Expiry time.Duration
}

type CryptoEnv struct {
	Key     []byte // 32 bytes from ENCRYPTION_KEY hex
	Version int    // ENCRYPTION_KEY_VERSION
}

type CORSEnv struct {
	Origins []string
}

type RateLimitEnv struct {
	Max    int
	Window time.Duration
}

func (v Vars) IsDevelopment() bool {
	return v.App.Env == "development"
}

// Load reads .env (walk-up), validates, and populates Envs.
func Load() error {
	loadEnvFiles()

	databaseURL, errDatabaseURL := requireString("DATABASE_URL")
	jwtSecret, errJWTSecret := requireString("JWT_SECRET")
	encKey, errEncKey := requireEncryptionKey("ENCRYPTION_KEY")

	jwtExpiry, errJWTExpiry := optionalDuration("JWT_EXPIRY", "24h")
	rateLimitWindow, errRateLimitWindow := optionalDuration("RATE_LIMIT_WINDOW", "1m")
	rateLimitMax, errRateLimitMax := optionalInt("RATE_LIMIT_MAX", "100")
	encKeyVersion, errEncKeyVersion := optionalInt("ENCRYPTION_KEY_VERSION", "1")

	origins, errCORS := parseCORSOrigins(os.Getenv("CORS_ORIGINS"))

	if err := joinErrors(errDatabaseURL, errJWTSecret, errEncKey, errJWTExpiry, errRateLimitWindow, errRateLimitMax, errEncKeyVersion, errCORS); err != nil {
		return err
	}

	if rateLimitMax < 1 {
		return fmt.Errorf("RATE_LIMIT_MAX must be >= 1")
	}
	if encKeyVersion < 1 {
		return fmt.Errorf("ENCRYPTION_KEY_VERSION must be >= 1")
	}

	Envs = Vars{
		App: AppEnv{
			Env:  optionalString("APP_ENV", "development"),
			Port: optionalString("PORT", "8080"),
		},
		Database: DatabaseEnv{
			URL: databaseURL,
		},
		JWT: JWTEnv{
			Secret: jwtSecret,
			Expiry: jwtExpiry,
		},
		Crypto: CryptoEnv{
			Key:     encKey,
			Version: encKeyVersion,
		},
		CORS: CORSEnv{
			Origins: origins,
		},
		RateLimit: RateLimitEnv{
			Max:    rateLimitMax,
			Window: rateLimitWindow,
		},
	}

	return nil
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

func requireString(key string) (string, error) {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return v, nil
}

func requireEncryptionKey(key string) ([]byte, error) {
	raw, err := requireString(key)
	if err != nil {
		return nil, err
	}
	decoded, err := hex.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("%s must be hex-encoded: %w", key, err)
	}
	if len(decoded) != 32 {
		return nil, fmt.Errorf("%s must decode to exactly 32 bytes (got %d)", key, len(decoded))
	}
	return decoded, nil
}

func optionalString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func optionalDuration(key, fallback string) (time.Duration, error) {
	raw := optionalString(key, fallback)
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}
	return d, nil
}

func optionalInt(key, fallback string) (int, error) {
	raw := optionalString(key, fallback)
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}
	return n, nil
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

func joinErrors(errs ...error) error {
	var nonNil []error
	for _, err := range errs {
		if err != nil {
			nonNil = append(nonNil, err)
		}
	}
	if len(nonNil) == 0 {
		return nil
	}
	return errors.Join(nonNil...)
}
