package env_test

import (
	"os"
	"strings"
	"testing"

	"github.com/envm-vault/envm-vault/api/internal/env"
)

const testEncryptionKey = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func clearAppEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"DATABASE_URL", "JWT_SECRET", "ENCRYPTION_KEY", "ENCRYPTION_KEY_VERSION",
		"APP_ENV", "PORT", "JWT_EXPIRY", "CORS_ORIGINS", "RATE_LIMIT_MAX", "RATE_LIMIT_WINDOW",
	} {
		os.Unsetenv(key)
	}
}

func TestLoad_FromRepoRootEnv(t *testing.T) {
	if os.Getenv("DATABASE_URL") == "" && os.Getenv("JWT_SECRET") == "" {
		wd, _ := os.Getwd()
		t.Skipf("DATABASE_URL/JWT_SECRET not in environment; run from repo with .env present (cwd=%s)", wd)
	}

	if err := env.Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if env.Envs.Database.URL == "" {
		t.Fatal("expected DATABASE_URL to be loaded")
	}
	if env.Envs.JWT.Secret == "" {
		t.Fatal("expected JWT_SECRET to be loaded")
	}
	if len(env.Envs.Crypto.Key) != 32 {
		t.Fatalf("Crypto.Key len = %d, want 32", len(env.Envs.Crypto.Key))
	}
}

func TestLoad_MissingEnv(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Chdir(t.TempDir())
	t.Cleanup(func() { t.Chdir(wd) })

	clearAppEnv(t)

	if err := env.Load(); err == nil {
		t.Fatal("expected error when required env is missing")
	} else {
		msg := err.Error()
		if !strings.Contains(msg, "DATABASE_URL is required") ||
			!strings.Contains(msg, "JWT_SECRET is required") ||
			!strings.Contains(msg, "ENCRYPTION_KEY is required") {
			t.Fatalf("error = %q, want DATABASE_URL, JWT_SECRET and ENCRYPTION_KEY required messages", msg)
		}
	}
}

func TestLoad_MissingOnlyJWTSecret(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Chdir(t.TempDir())
	t.Cleanup(func() { t.Chdir(wd) })

	clearAppEnv(t)
	t.Setenv("DATABASE_URL", "postgres://envm:envm@localhost:5432/envm?sslmode=disable")
	t.Setenv("ENCRYPTION_KEY", testEncryptionKey)

	if err := env.Load(); err == nil {
		t.Fatal("expected error when JWT_SECRET is missing")
	} else if !strings.Contains(err.Error(), "JWT_SECRET is required") {
		t.Fatalf("error = %q, want JWT_SECRET is required", err.Error())
	}
}

func TestLoad_InvalidEncryptionKey(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Chdir(t.TempDir())
	t.Cleanup(func() { t.Chdir(wd) })

	clearAppEnv(t)
	t.Setenv("DATABASE_URL", "postgres://envm:envm@localhost:5432/envm?sslmode=disable")
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("ENCRYPTION_KEY", "not-hex")

	if err := env.Load(); err == nil {
		t.Fatal("expected error for invalid ENCRYPTION_KEY")
	}
}

func TestLoad_CORSOriginsJSON(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Chdir(t.TempDir())
	t.Cleanup(func() { t.Chdir(wd) })

	t.Setenv("DATABASE_URL", "postgres://envm:envm@localhost:5432/envm?sslmode=disable")
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("ENCRYPTION_KEY", testEncryptionKey)
	t.Setenv("CORS_ORIGINS", `["http://localhost:3100","http://127.0.0.1:3100"]`)

	if err := env.Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(env.Envs.CORS.Origins) != 2 {
		t.Fatalf("CORS.Origins len = %d, want 2", len(env.Envs.CORS.Origins))
	}
	if env.Envs.CORS.Origins[0] != "http://localhost:3100" || env.Envs.CORS.Origins[1] != "http://127.0.0.1:3100" {
		t.Errorf("CORS.Origins = %#v", env.Envs.CORS.Origins)
	}
}

func TestLoad_CORSOriginsDefault(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Chdir(t.TempDir())
	t.Cleanup(func() { t.Chdir(wd) })

	t.Setenv("DATABASE_URL", "postgres://envm:envm@localhost:5432/envm?sslmode=disable")
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("ENCRYPTION_KEY", testEncryptionKey)
	t.Setenv("CORS_ORIGINS", "")

	if err := env.Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(env.Envs.CORS.Origins) != 1 || env.Envs.CORS.Origins[0] != "http://localhost:3100" {
		t.Errorf("CORS.Origins = %#v, want default localhost:3100", env.Envs.CORS.Origins)
	}
	if env.Envs.Crypto.Version != 1 {
		t.Errorf("Crypto.Version = %d, want 1", env.Envs.Crypto.Version)
	}
}

func TestLoad_CORSOriginsInvalidJSON(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Chdir(t.TempDir())
	t.Cleanup(func() { t.Chdir(wd) })

	t.Setenv("DATABASE_URL", "postgres://envm:envm@localhost:5432/envm?sslmode=disable")
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("ENCRYPTION_KEY", testEncryptionKey)
	t.Setenv("CORS_ORIGINS", "not-json")

	if err := env.Load(); err == nil {
		t.Fatal("expected error for invalid CORS_ORIGINS JSON")
	}
}

func TestLoad_RateLimitDefaults(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Chdir(t.TempDir())
	t.Cleanup(func() { t.Chdir(wd) })

	t.Setenv("DATABASE_URL", "postgres://envm:envm@localhost:5432/envm?sslmode=disable")
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("ENCRYPTION_KEY", testEncryptionKey)
	t.Setenv("CORS_ORIGINS", "")
	os.Unsetenv("RATE_LIMIT_MAX")
	os.Unsetenv("RATE_LIMIT_WINDOW")

	if err := env.Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if env.Envs.RateLimit.Max != 100 {
		t.Errorf("RateLimit.Max = %d, want 100", env.Envs.RateLimit.Max)
	}
	if env.Envs.RateLimit.Window.String() != "1m0s" {
		t.Errorf("RateLimit.Window = %v, want 1m", env.Envs.RateLimit.Window)
	}
}
