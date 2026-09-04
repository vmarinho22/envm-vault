package config

import (
	"os"
	"testing"
)

func TestLoad_FromRepoRootEnv(t *testing.T) {
	if os.Getenv("DATABASE_URL") == "" && os.Getenv("JWT_SECRET") == "" {
		// Walk-up loader should find .env at monorepo root when running from api/
		wd, _ := os.Getwd()
		t.Skipf("DATABASE_URL/JWT_SECRET not in environment; run from repo with .env present (cwd=%s)", wd)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.DatabaseURL == "" {
		t.Fatal("expected DATABASE_URL to be loaded")
	}
	if cfg.JWTSecret == "" {
		t.Fatal("expected JWT_SECRET to be loaded")
	}
}

func TestLoad_MissingEnv(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Chdir(t.TempDir())

	for _, key := range []string{"DATABASE_URL", "JWT_SECRET", "APP_ENV", "PORT", "JWT_EXPIRY", "CORS_ORIGINS"} {
		os.Unsetenv(key)
	}

	_, err = Load()
	if err == nil {
		t.Fatal("expected error when DATABASE_URL is missing")
	}

	t.Chdir(wd)
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
	t.Setenv("CORS_ORIGINS", `["http://localhost:3100","http://127.0.0.1:3100"]`)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(cfg.CORSOrigins) != 2 {
		t.Fatalf("CORSOrigins len = %d, want 2", len(cfg.CORSOrigins))
	}
	if cfg.CORSOrigins[0] != "http://localhost:3100" || cfg.CORSOrigins[1] != "http://127.0.0.1:3100" {
		t.Errorf("CORSOrigins = %#v", cfg.CORSOrigins)
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
	t.Setenv("CORS_ORIGINS", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(cfg.CORSOrigins) != 1 || cfg.CORSOrigins[0] != "http://localhost:3100" {
		t.Errorf("CORSOrigins = %#v, want default localhost:3100", cfg.CORSOrigins)
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
	t.Setenv("CORS_ORIGINS", "not-json")

	_, err = Load()
	if err == nil {
		t.Fatal("expected error for invalid CORS_ORIGINS JSON")
	}
}
