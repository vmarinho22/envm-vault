package tag

import (
	"github.com/envm-vault/envm-vault/api/internal/handler/httperr"
	"github.com/envm-vault/envm-vault/api/internal/middleware"
	tagservice "github.com/envm-vault/envm-vault/api/internal/service/tag"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	tag *tagservice.Service
}

func New(tag *tagservice.Service) *Handler {
	return &Handler{tag: tag}
}

type createRequest struct {
	Title string `json:"title"`
	Color string `json:"color"`
}

type updateRequest struct {
	Title *string `json:"title"`
	Color *string `json:"color"`
}

func (h *Handler) List(c *fiber.Ctx) error {
	title := c.Query("title")
	tags, err := h.tag.List(c.UserContext(), title)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.JSON(tags)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	tag, err := h.tag.Create(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		req.Title,
		req.Color,
	)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(tag)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	var req updateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	tag, err := h.tag.Update(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		c.Params("id"),
		req.Title,
		req.Color,
	)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.JSON(tag)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	if err := h.tag.Delete(
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
