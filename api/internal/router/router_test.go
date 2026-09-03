package router

import (
	"testing"

	authhandler "github.com/envm-vault/envm-vault/api/internal/handler/auth"
	healthhandler "github.com/envm-vault/envm-vault/api/internal/handler/health"
	userhandler "github.com/envm-vault/envm-vault/api/internal/handler/user"
	"github.com/gofiber/fiber/v2"
)

func TestRegister_RoutesRegistered(t *testing.T) {
	app := fiber.New()
	Register(app, Deps{
		Health: &healthhandler.Handler{},
		Auth:   &authhandler.Handler{},
		User:   &userhandler.Handler{},
	})

	want := map[string]bool{
		"GET /v1/health":         false,
		"GET /health":            false,
		"POST /v1/auth/login":    false,
		"POST /v1/auth/verify":   false,
		"POST /v1/auth/register": false,
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
