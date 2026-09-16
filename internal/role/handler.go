package role

import "github.com/gofiber/fiber/v2"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateRole(c *fiber.Ctx) error {
	var input CreateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	role, err := h.service.Create(c.Context(), &input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(role)
}

func (h *Handler) GetRoles(c *fiber.Ctx) error {
	roles, err := h.service.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	responses := make([]Response, 0, len(roles))
	for _, r := range roles {
		if r == nil {
			continue
		}
		responses = append(responses, Response{ID: r.ID, Name: r.Name})
	}
	return c.Status(fiber.StatusOK).JSON(responses)
}
