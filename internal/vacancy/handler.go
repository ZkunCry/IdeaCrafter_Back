package vacancy

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func actorID(c *fiber.Ctx) (uint, bool) {
	userID, ok := c.Locals("user_id").(uint)
	return userID, ok
}

func respondError(c *fiber.Ctx, err error) error {
	if errors.Is(err, ErrForbidden) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
}

func (h *Handler) CreateVacancy(c *fiber.Ctx) error {
	userID, ok := actorID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	var input CreateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	vacancy, err := h.service.Create(c.Context(), &input, userID)
	if err != nil {
		return respondError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(NewResponse(vacancy))
}

func (h *Handler) GetVacancyByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	vacancy, err := h.service.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(NewResponse(vacancy))
}

func (h *Handler) GetVacanciesByStartup(c *fiber.Ctx) error {
	startupID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid startup id"})
	}
	vacancies, err := h.service.GetByStartupID(c.Context(), uint(startupID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(NewResponses(vacancies))
}

func (h *Handler) UpdateVacancy(c *fiber.Ctx) error {
	userID, ok := actorID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	var input UpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	vacancy, err := h.service.Update(c.Context(), uint(id), &input, userID)
	if err != nil {
		return respondError(c, err)
	}
	return c.JSON(NewResponse(vacancy))
}

func (h *Handler) DeleteVacancy(c *fiber.Ctx) error {
	userID, ok := actorID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := h.service.Delete(c.Context(), uint(id), userID); err != nil {
		return respondError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "vacancy deleted"})
}
