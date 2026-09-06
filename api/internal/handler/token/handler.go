package token

import (
	"github.com/envm-vault/envm-vault/api/internal/handler/httperr"
	"github.com/envm-vault/envm-vault/api/internal/middleware"
	tokenservice "github.com/envm-vault/envm-vault/api/internal/service/token"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	tokens *tokenservice.Service
}

func New(tokens *tokenservice.Service) *Handler {
	return &Handler{tokens: tokens}
}

func (h *Handler) GetMine(c *fiber.Ctx) error {
	resp, err := h.tokens.GetMine(c.UserContext(), middleware.GetUserID(c))
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.JSON(resp)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	resp, err := h.tokens.Create(c.UserContext(), middleware.GetUserID(c))
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *Handler) RevokeOwn(c *fiber.Ctx) error {
	if err := h.tokens.RevokeOwn(c.UserContext(), middleware.GetUserID(c)); err != nil {
		return httperr.Write(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) RevokeForUser(c *fiber.Ctx) error {
	err := h.tokens.RevokeForUser(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserRole(c),
		c.Params("id"),
	)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
