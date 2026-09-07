package router_test

import (
	"net/http/httptest"
	"testing"

	audithandler "github.com/envm-vault/envm-vault/api/internal/handler/audit"
	authhandler "github.com/envm-vault/envm-vault/api/internal/handler/auth"
	environmenthandler "github.com/envm-vault/envm-vault/api/internal/handler/environment"
	healthhandler "github.com/envm-vault/envm-vault/api/internal/handler/health"
	taghandler "github.com/envm-vault/envm-vault/api/internal/handler/tag"
	teamhandler "github.com/envm-vault/envm-vault/api/internal/handler/team"
	tokenhandler "github.com/envm-vault/envm-vault/api/internal/handler/token"
	userhandler "github.com/envm-vault/envm-vault/api/internal/handler/user"
	"github.com/envm-vault/envm-vault/api/internal/middleware"
	"github.com/envm-vault/envm-vault/api/internal/router"
	"github.com/gofiber/fiber/v2"
)

func TestRegister_RoutesRegistered(t *testing.T) {
	app := fiber.New()
	router.Register(app, router.Deps{
		Health:      &healthhandler.Handler{},
		Auth:        &authhandler.Handler{},
		User:        &userhandler.Handler{},
		Token:       &tokenhandler.Handler{},
		Audit:       &audithandler.Handler{},
		Team:        &teamhandler.Handler{},
		Tag:         &taghandler.Handler{},
		Environment: &environmenthandler.Handler{},
	})

	want := map[string]bool{
		"GET /v1/health":              false,
		"GET /health":                 false,
		"POST /v1/auth/login":         false,
		"POST /v1/auth/verify":        false,
		"POST /v1/auth/register":      false,
		"GET /v1/users":               false,
		"DELETE /v1/users/:id/token":  false,
		"POST /v1/users/:id/block":    false,
		"POST /v1/users/:id/unblock":  false,
		"GET /v1/tokens/me":           false,
		"POST /v1/tokens":             false,
		"DELETE /v1/tokens":           false,
		"GET /v1/audit-logs":          false,
		"GET /v1/teams":               false,
		"POST /v1/teams":              false,
		"PATCH /v1/teams/:id":         false,
		"DELETE /v1/teams/:id":        false,
		"GET /v1/tags":                false,
		"POST /v1/tags":               false,
		"PATCH /v1/tags/:id":          false,
		"DELETE /v1/tags/:id":         false,
		"GET /v1/environments":        false,
		"POST /v1/environments":       false,
		"PATCH /v1/environments/:id":  false,
		"DELETE /v1/environments/:id": false,
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
	router.Register(app, router.Deps{
		Health:      &healthhandler.Handler{},
		Auth:        &authhandler.Handler{},
		User:        &userhandler.Handler{},
		Token:       &tokenhandler.Handler{},
		Audit:       &audithandler.Handler{},
		Team:        &teamhandler.Handler{},
		Tag:         &taghandler.Handler{},
		Environment: &environmenthandler.Handler{},
		JWT:         middleware.NewJWTMiddleware("test-secret", nil, nil),
	})

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
