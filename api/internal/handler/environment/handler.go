package environment

import (
	"github.com/envm-vault/envm-vault/api/internal/handler/httperr"
	"github.com/envm-vault/envm-vault/api/internal/middleware"
	environmentservice "github.com/envm-vault/envm-vault/api/internal/service/environment"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	environment *environmentservice.Service
}

func New(environment *environmentservice.Service) *Handler {
	return &Handler{environment: environment}
}

type createRequest struct {
	Title string `json:"title"`
}

type updateRequest struct {
	Title *string `json:"title"`
}

func (h *Handler) List(c *fiber.Ctx) error {
	title := c.Query("title")
	environments, err := h.environment.List(c.UserContext(), title)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.JSON(environments)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	environment, err := h.environment.Create(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		req.Title,
	)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(environment)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	var req updateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	environment, err := h.environment.Update(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		c.Params("id"),
		req.Title,
	)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.JSON(environment)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	if err := h.environment.Delete(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		c.Params("id"),
	); err != nil {
		return httperr.Write(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
