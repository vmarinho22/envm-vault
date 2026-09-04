package env

import (
	"os"
	"strings"
	"testing"
)

func TestLoad_FromRepoRootEnv(t *testing.T) {
	if os.Getenv("DATABASE_URL") == "" && os.Getenv("JWT_SECRET") == "" {
		wd, _ := os.Getwd()
		t.Skipf("DATABASE_URL/JWT_SECRET not in environment; run from repo with .env present (cwd=%s)", wd)
	}

	if err := Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if Envs.Database.URL == "" {
		t.Fatal("expected DATABASE_URL to be loaded")
	}
	if Envs.JWT.Secret == "" {
		t.Fatal("expected JWT_SECRET to be loaded")
	}
}

func TestLoad_MissingEnv(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Chdir(t.TempDir())
	t.Cleanup(func() { t.Chdir(wd) })

	for _, key := range []string{"DATABASE_URL", "JWT_SECRET", "APP_ENV", "PORT", "JWT_EXPIRY", "CORS_ORIGINS", "RATE_LIMIT_MAX", "RATE_LIMIT_WINDOW"} {
		os.Unsetenv(key)
	}

	if err := Load(); err == nil {
		t.Fatal("expected error when required env is missing")
	} else {
		msg := err.Error()
		if !strings.Contains(msg, "DATABASE_URL is required") || !strings.Contains(msg, "JWT_SECRET is required") {
			t.Fatalf("error = %q, want both DATABASE_URL and JWT_SECRET required messages", msg)
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

	for _, key := range []string{"DATABASE_URL", "JWT_SECRET", "APP_ENV", "PORT", "JWT_EXPIRY", "CORS_ORIGINS", "RATE_LIMIT_MAX", "RATE_LIMIT_WINDOW"} {
		os.Unsetenv(key)
	}
	t.Setenv("DATABASE_URL", "postgres://envm:envm@localhost:5432/envm?sslmode=disable")

	if err := Load(); err == nil {
		t.Fatal("expected error when JWT_SECRET is missing")
	} else if !strings.Contains(err.Error(), "JWT_SECRET is required") {
		t.Fatalf("error = %q, want JWT_SECRET is required", err.Error())
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
	t.Setenv("CORS_ORIGINS", `["http://localhost:3100","http://127.0.0.1:3100"]`)

	if err := Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(Envs.CORS.Origins) != 2 {
		t.Fatalf("CORS.Origins len = %d, want 2", len(Envs.CORS.Origins))
	}
	if Envs.CORS.Origins[0] != "http://localhost:3100" || Envs.CORS.Origins[1] != "http://127.0.0.1:3100" {
		t.Errorf("CORS.Origins = %#v", Envs.CORS.Origins)
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

	if err := Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(Envs.CORS.Origins) != 1 || Envs.CORS.Origins[0] != "http://localhost:3100" {
		t.Errorf("CORS.Origins = %#v, want default localhost:3100", Envs.CORS.Origins)
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

	if err := Load(); err == nil {
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
	t.Setenv("CORS_ORIGINS", "")
	os.Unsetenv("RATE_LIMIT_MAX")
	os.Unsetenv("RATE_LIMIT_WINDOW")

	if err := Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if Envs.RateLimit.Max != 100 {
		t.Errorf("RateLimit.Max = %d, want 100", Envs.RateLimit.Max)
	}
	if Envs.RateLimit.Window.String() != "1m0s" {
		t.Errorf("RateLimit.Window = %v, want 1m", Envs.RateLimit.Window)
	}
}
