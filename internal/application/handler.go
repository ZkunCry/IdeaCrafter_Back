package application

import (
	"errors"
	"strconv"

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

func (h *Handler) CreateApplication(c *fiber.Ctx) error {
	var input CreateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	userID, ok := actorID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	input.UserID = userID
	application, err := h.service.Create(c.Context(), &input)
	if err != nil {
		return respondError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(NewResponse(application))
}

func (h *Handler) UpdateApplication(c *fiber.Ctx) error {
	userID, ok := actorID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	var input UpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	application, err := h.service.Update(c.Context(), uint(id), &input, userID)
	if err != nil {
		return respondError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(NewResponse(application))
}

func (h *Handler) UpdateApplicationStatus(c *fiber.Ctx) error {
	userID, ok := actorID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	var input UpdateStatusInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	application, err := h.service.UpdateStatus(c.Context(), uint(id), &input, userID)
	if err != nil {
		return respondError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(NewResponse(application))
}

func (h *Handler) GetApplicationByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	application, err := h.service.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(NewResponse(application))
}


func (h *Handler) GetMyApplications(c *fiber.Ctx) error {
	userID, ok := actorID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	applications, err := h.service.GetByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	briefs, err := h.service.GetStartupBriefs(c.Context(), applications)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	responses := NewResponses(applications)
	for i := range responses {
		if brief, ok := briefs[responses[i].StartupID]; ok {
			responses[i].Startup = &brief
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"items": responses})
}

func (h *Handler) GetStartupApplications(c *fiber.Ctx) error {
	userID, ok := actorID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	startupID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid startup id"})
	}
	applications, err := h.service.GetByStartupID(c.Context(), uint(startupID), userID)
	if err != nil {
		return respondError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"items": NewResponses(applications)})
}

func (h *Handler) DeleteApplication(c *fiber.Ctx) error {
	userID, ok := actorID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := h.service.Delete(c.Context(), uint(id), userID); err != nil {
		return respondError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "application withdrawn"})
}
