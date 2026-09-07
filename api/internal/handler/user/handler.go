package user

import (
	"github.com/envm-vault/envm-vault/api/internal/handler/httperr"
	"github.com/envm-vault/envm-vault/api/internal/middleware"
	userservice "github.com/envm-vault/envm-vault/api/internal/service/user"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	user *userservice.Service
}

func New(user *userservice.Service) *Handler {
	return &Handler{user: user}
}

type registerRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	role := middleware.GetUserRole(c)
	user, err := h.user.Register(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		role,
		req.Email,
		req.Role,
	)
	if err != nil {
		return httperr.Write(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *Handler) List(c *fiber.Ctx) error {
	users, err := h.user.List(c.UserContext(), middleware.GetUserRole(c))
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.JSON(users)
}

func (h *Handler) Block(c *fiber.Ctx) error {
	user, err := h.user.Block(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		c.Params("id"),
	)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.JSON(user)
}

func (h *Handler) Unblock(c *fiber.Ctx) error {
	user, err := h.user.Unblock(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		c.Params("id"),
	)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.JSON(user)
}
