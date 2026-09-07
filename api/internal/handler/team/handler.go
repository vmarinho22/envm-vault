package team

import (
	"github.com/envm-vault/envm-vault/api/internal/handler/httperr"
	"github.com/envm-vault/envm-vault/api/internal/middleware"
	teamservice "github.com/envm-vault/envm-vault/api/internal/service/team"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	team *teamservice.Service
}

func New(team *teamservice.Service) *Handler {
	return &Handler{team: team}
}

type createRequest struct {
	Name          string `json:"name"`
	AccountableID string `json:"accountable_id"`
}

type updateRequest struct {
	Name          *string `json:"name"`
	AccountableID *string `json:"accountable_id"`
}

func (h *Handler) List(c *fiber.Ctx) error {
	name := c.Query("name")
	teams, err := h.team.List(c.UserContext(), name)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.JSON(teams)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	team, err := h.team.Create(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		req.Name,
		req.AccountableID,
	)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(team)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	var req updateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	team, err := h.team.Update(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		c.Params("id"),
		req.Name,
		req.AccountableID,
	)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.JSON(team)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	if err := h.team.Delete(
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
