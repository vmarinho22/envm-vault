package router_test

import (
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	audithandler "github.com/envm-vault/envm-vault/api/internal/handler/audit"
	authhandler "github.com/envm-vault/envm-vault/api/internal/handler/auth"
	environmenthandler "github.com/envm-vault/envm-vault/api/internal/handler/environment"
	healthhandler "github.com/envm-vault/envm-vault/api/internal/handler/health"
	projecthandler "github.com/envm-vault/envm-vault/api/internal/handler/project"
	secrethandler "github.com/envm-vault/envm-vault/api/internal/handler/secret"
	taghandler "github.com/envm-vault/envm-vault/api/internal/handler/tag"
	teamhandler "github.com/envm-vault/envm-vault/api/internal/handler/team"
	tokenhandler "github.com/envm-vault/envm-vault/api/internal/handler/token"
	userhandler "github.com/envm-vault/envm-vault/api/internal/handler/user"
	"github.com/envm-vault/envm-vault/api/internal/middleware"
	"github.com/envm-vault/envm-vault/api/internal/router"
	"github.com/gofiber/fiber/v2"
)

func emptyDeps(jwt *middleware.JWTMiddleware, apiKey *middleware.APIKeyMiddleware) router.Deps {
	return router.Deps{
		Health:      &healthhandler.Handler{},
		Auth:        &authhandler.Handler{},
		User:        &userhandler.Handler{},
		Token:       &tokenhandler.Handler{},
		Audit:       &audithandler.Handler{},
		Team:        &teamhandler.Handler{},
		Tag:         &taghandler.Handler{},
		Environment: &environmenthandler.Handler{},
		Project:     &projecthandler.Handler{},
		Secret:      &secrethandler.Handler{},
		JWT:         jwt,
		APIKey:      apiKey,
	}
}

func TestRegister_RoutesRegistered(t *testing.T) {
	app := fiber.New()
	router.Register(app, emptyDeps(
		middleware.NewJWTMiddleware("test-secret", nil, nil),
		middleware.NewAPIKeyMiddleware(nil, nil, nil),
	))

	want := map[string]bool{
		"GET /v1/health":                 false,
		"GET /health":                    false,
		"POST /v1/auth/login":            false,
		"POST /v1/auth/verify":           false,
		"POST /v1/auth/register":         false,
		"GET /v1/users":                  false,
		"DELETE /v1/users/:id/token":     false,
		"POST /v1/users/:id/block":       false,
		"POST /v1/users/:id/unblock":     false,
		"GET /v1/tokens/me":              false,
		"POST /v1/tokens":                false,
		"DELETE /v1/tokens":              false,
		"GET /v1/audit-logs":             false,
		"GET /v1/teams":                  false,
		"POST /v1/teams":                 false,
		"PATCH /v1/teams/:id":            false,
		"DELETE /v1/teams/:id":           false,
		"GET /v1/tags":                   false,
		"POST /v1/tags":                  false,
		"PATCH /v1/tags/:id":             false,
		"DELETE /v1/tags/:id":            false,
		"GET /v1/environments":           false,
		"POST /v1/environments":          false,
		"PATCH /v1/environments/:id":     false,
		"DELETE /v1/environments/:id":    false,
		"GET /v1/projects":               false,
		"GET /v1/projects/by-slug/:slug": false,
		"POST /v1/projects":              false,
		"PATCH /v1/projects/:id":         false,
		"DELETE /v1/projects/:id":        false,
		"GET /v1/projects/:id/secrets":   false,
		"POST /v1/projects/:id/secrets":  false,
		"GET /v1/projects/:slug/env":     false,
		"PUT /v1/projects/:id/secrets/:secretId/prefixes/:environmentId": false,
	}

	for _, route := range app.GetRoutes() {
		key := route.Method + " " + route.Path
		if _, ok := want[key]; ok {
			want[key] = true
		}
	}

	for path, found := range want {
		if !found {
			t.Errorf("route not registered: %s", path)
		}
	}
}

func TestRegister_TokensMeRequiresAuth(t *testing.T) {
	app := fiber.New()
	router.Register(app, emptyDeps(
		middleware.NewJWTMiddleware("test-secret", nil, nil),
		middleware.NewAPIKeyMiddleware(nil, nil, nil),
	))

	req := httptest.NewRequest(fiber.MethodGet, "/v1/tokens/me", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("status = %d, want %d (route missing would be 404)", resp.StatusCode, fiber.StatusUnauthorized)
	}

	postReq := httptest.NewRequest(fiber.MethodPost, "/v1/tokens", nil)
	postResp, err := app.Test(postReq, -1)
	if err != nil {
		t.Fatalf("post request: %v", err)
	}
	defer postResp.Body.Close()

	if postResp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("POST status = %d, want %d", postResp.StatusCode, fiber.StatusUnauthorized)
	}
}

func TestRegister_DownloadRequiresAPIKey(t *testing.T) {
	app := fiber.New()
	router.Register(app, emptyDeps(
		middleware.NewJWTMiddleware("test-secret", nil, nil),
		middleware.NewAPIKeyMiddleware(nil, nil, nil),
	))

	// Non-ev_ bearer must hit APIKey (prefix check), not JWT parse.
	req := httptest.NewRequest(fiber.MethodGet, "/v1/projects/my-app/env?environment=dev", nil)
	req.Header.Set("Authorization", "Bearer not-an-api-key")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusUnauthorized)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	got := string(body)
	if !strings.Contains(got, "invalid or revoked access token") {
		t.Fatalf("body = %q, want API-key error (not JWT)", got)
	}
	if strings.Contains(got, "invalid or expired token") {
		t.Fatalf("download hit JWT middleware: %q", got)
	}
}

func TestMountFrontend_ServesIndexAndKeepsAPI(t *testing.T) {
	dir := t.TempDir()
	indexPath := dir + "/index.html"
	if err := os.WriteFile(indexPath, []byte("<!doctype html><title>ui</title>"), 0o644); err != nil {
		t.Fatalf("write index: %v", err)
	}

	app := fiber.New()
	router.Register(app, emptyDeps(
		middleware.NewJWTMiddleware("test-secret", nil, nil),
		middleware.NewAPIKeyMiddleware(nil, nil, nil),
	))
	router.MountFrontend(app, dir)

	spaReq := httptest.NewRequest(fiber.MethodGet, "/projects/demo", nil)
	spaResp, err := app.Test(spaReq, -1)
	if err != nil {
		t.Fatalf("spa request: %v", err)
	}
	defer spaResp.Body.Close()
	if spaResp.StatusCode != fiber.StatusOK {
		t.Fatalf("spa status = %d, want %d", spaResp.StatusCode, fiber.StatusOK)
	}
	body, _ := io.ReadAll(spaResp.Body)
	if !strings.Contains(string(body), "<title>ui</title>") {
		t.Fatalf("spa body = %q, want index html", body)
	}

	apiReq := httptest.NewRequest(fiber.MethodGet, "/v1/tokens/me", nil)
	apiResp, err := app.Test(apiReq, -1)
	if err != nil {
		t.Fatalf("api request: %v", err)
	}
	defer apiResp.Body.Close()
	if apiResp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("api status = %d, want %d (route must not be SPA)", apiResp.StatusCode, fiber.StatusUnauthorized)
	}
	apiBody, _ := io.ReadAll(apiResp.Body)
	if strings.Contains(string(apiBody), "<title>ui</title>") {
		t.Fatalf("API route was shadowed by SPA")
	}
}

func TestMountFrontend_EmptyDistIsNoop(t *testing.T) {
	app := fiber.New()
	router.MountFrontend(app, "")
	router.MountFrontend(app, "/path/does/not/exist")

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("status = %d, want 404 when frontend is not mounted", resp.StatusCode)
	}
}
