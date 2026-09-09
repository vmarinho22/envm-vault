package httperr_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/vmarinho22/envm-vault/api/internal/handler/httperr"
	"github.com/vmarinho22/envm-vault/api/internal/service"
	"github.com/gofiber/fiber/v2"
)

func TestWrite_TooManyAttempts(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return httperr.Write(c, service.ErrTooManyAttempts)
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusTooManyRequests)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	var payload map[string]string
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	want := "too many verification attempts; request a new code"
	if payload["error"] != want {
		t.Errorf("error = %q, want %q", payload["error"], want)
	}
}

func TestWrite_TokenNotFound(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return httperr.Write(c, service.ErrTokenNotFound)
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusNotFound)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	var payload map[string]string
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	want := "access token not found"
	if payload["error"] != want {
		t.Errorf("error = %q, want %q", payload["error"], want)
	}
}

func TestAppHandler_NotFoundKeepsStatus(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: httperr.AppHandler})

	req := httptest.NewRequest(fiber.MethodGet, "/missing", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusNotFound)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	var payload map[string]string
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if payload["error"] == "" || payload["error"] == "internal server error" {
		t.Errorf("error = %q, want Fiber not-found message", payload["error"])
	}
}

func TestAppHandler_UnexpectedErrorIsGeneric500(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: httperr.AppHandler})
	app.Get("/", func(c *fiber.Ctx) error {
		return errors.New("boom")
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusInternalServerError)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	var payload map[string]string
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if payload["error"] != "internal server error" {
		t.Errorf("error = %q, want generic internal server error", payload["error"])
	}
}
