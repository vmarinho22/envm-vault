package router_test

import (
	"net/http/httptest"
	"testing"

	authhandler "github.com/envm-vault/envm-vault/api/internal/handler/auth"
	healthhandler "github.com/envm-vault/envm-vault/api/internal/handler/health"
	tokenhandler "github.com/envm-vault/envm-vault/api/internal/handler/token"
	userhandler "github.com/envm-vault/envm-vault/api/internal/handler/user"
	"github.com/envm-vault/envm-vault/api/internal/middleware"
	"github.com/envm-vault/envm-vault/api/internal/router"
	"github.com/gofiber/fiber/v2"
)

func TestRegister_RoutesRegistered(t *testing.T) {
	app := fiber.New()
	router.Register(app, router.Deps{
		Health: &healthhandler.Handler{},
		Auth:   &authhandler.Handler{},
		User:   &userhandler.Handler{},
		Token:  &tokenhandler.Handler{},
	})

	want := map[string]bool{
		"GET /v1/health":             false,
		"GET /health":                false,
		"POST /v1/auth/login":        false,
		"POST /v1/auth/verify":       false,
		"POST /v1/auth/register":     false,
		"GET /v1/users":              false,
		"DELETE /v1/users/:id/token": false,
		"GET /v1/tokens/me":          false,
		"POST /v1/tokens":            false,
		"DELETE /v1/tokens":          false,
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
		Health: &healthhandler.Handler{},
		Auth:   &authhandler.Handler{},
		User:   &userhandler.Handler{},
		Token:  &tokenhandler.Handler{},
		JWT:    middleware.NewJWTMiddleware("test-secret"),
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
		t.Fatalf("POST /v1/tokens status = %d, want %d (route missing would be 404)", postResp.StatusCode, fiber.StatusUnauthorized)
	}
}
