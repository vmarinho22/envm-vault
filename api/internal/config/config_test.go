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

	for _, key := range []string{"DATABASE_URL", "JWT_SECRET", "APP_ENV", "PORT", "JWT_EXPIRY"} {
		os.Unsetenv(key)
	}

	_, err = Load()
	if err == nil {
		t.Fatal("expected error when DATABASE_URL is missing")
	}

	t.Chdir(wd)
}
