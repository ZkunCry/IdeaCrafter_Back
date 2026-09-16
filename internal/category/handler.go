package category

import (
	"math"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateCategory(c *fiber.Ctx) error {
	var input CreateInput
	if err := c.BodyParser(&input); err != nil {
		return err
	}
	category, err := h.service.Create(c.Context(), input)
	if err != nil {
		return err
	}
	return c.JSON(category)
}

func (h *Handler) GetCategory(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid category id"})
	}
	category, err := h.service.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(category)
}

func (h *Handler) GetAllCategories(c *fiber.Ctx) error {
	input := ListInput{Limit: 10, Offset: 0}
	if err := c.QueryParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	categories, totalCount, err := h.service.GetAll(c.Context(), input.SearchString, input.Limit, input.Offset)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(input.Limit)))
	response := ListResponse{Items: categories, Total: totalPages}
	return c.Status(fiber.StatusOK).JSON(response)
}
