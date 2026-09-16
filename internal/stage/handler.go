package stage

import "github.com/gofiber/fiber/v2"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateStage(c *fiber.Ctx) error {
	var input CreateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	stage, err := h.service.Create(c.Context(), &input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(stage)
}

func (h *Handler) GetList(c *fiber.Ctx) error {
	stages, err := h.service.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// Serialise through Response so the list uses the same lowercase "id" the
	// startup payload exposes instead of the raw gorm.Model "ID".
	responses := make([]Response, 0, len(stages))
	for _, stage := range stages {
		if stage == nil {
			continue
		}
		responses = append(responses, Response{ID: stage.ID, Name: stage.Name})
	}
	return c.Status(fiber.StatusOK).JSON(responses)
}
