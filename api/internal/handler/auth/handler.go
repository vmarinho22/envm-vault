package auth

import (
	"github.com/envm-vault/envm-vault/api/internal/handler/httperr"
	authservice "github.com/envm-vault/envm-vault/api/internal/service/auth"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	auth *authservice.Service
}

func New(auth *authservice.Service) *Handler {
	return &Handler{auth: auth}
}

type loginRequest struct {
	Email string `json:"email"`
}

type verifyRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.auth.RequestLogin(c.Context(), req.Email); err != nil {
		return httperr.Write(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "verification code sent",
	})
}

func (h *Handler) Verify(c *fiber.Ctx) error {
	var req verifyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	resp, err := h.auth.VerifyLogin(c.Context(), req.Email, req.Code)
	if err != nil {
		return httperr.Write(c, err)
	}

	return c.JSON(resp)
}
