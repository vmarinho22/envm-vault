package middleware_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/envm-vault/envm-vault/api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func TestRateLimiter_BlocksAfterMax(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.NewRateLimiter(2, time.Minute))
	app.Get("/v1/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	app.Get("/v1/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(fiber.MethodGet, "/v1/ping", nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("request %d: %v", i+1, err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("request %d status = %d, want %d", i+1, resp.StatusCode, fiber.StatusOK)
		}
		_ = resp.Body.Close()
	}

	req := httptest.NewRequest(fiber.MethodGet, "/v1/ping", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request 3: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusTooManyRequests {
		t.Fatalf("request 3 status = %d, want %d", resp.StatusCode, fiber.StatusTooManyRequests)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if !strings.Contains(string(body), "rate limit exceeded") {
		t.Errorf("body = %q, want rate limit exceeded", body)
	}

	healthReq := httptest.NewRequest(fiber.MethodGet, "/health", nil)
	healthResp, err := app.Test(healthReq, -1)
	if err != nil {
		t.Fatalf("health: %v", err)
	}
	defer healthResp.Body.Close()
	if healthResp.StatusCode != fiber.StatusOK {
		t.Fatalf("health status = %d, want %d", healthResp.StatusCode, fiber.StatusOK)
	}

	v1HealthReq := httptest.NewRequest(fiber.MethodGet, "/v1/health", nil)
	v1HealthResp, err := app.Test(v1HealthReq, -1)
	if err != nil {
		t.Fatalf("v1 health: %v", err)
	}
	defer v1HealthResp.Body.Close()
	if v1HealthResp.StatusCode != fiber.StatusOK {
		t.Fatalf("v1 health status = %d, want %d", v1HealthResp.StatusCode, fiber.StatusOK)
	}
}
