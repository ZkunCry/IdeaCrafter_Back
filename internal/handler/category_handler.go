package handler

import (
	"math"
	"startup_back/internal/dto"
	"startup_back/internal/service"

	"github.com/gofiber/fiber/v2"
)

type CategorynHandler struct {
	services *service.Services
}

func NewCategorynHandler(services *service.Services) *CategorynHandler {
	return &CategorynHandler{services: services}
}

func (h *CategorynHandler) CreateCategory(c *fiber.Ctx) error {
	var input dto.CreateCategoryInput
	if err := c.BodyParser(&input); err != nil {
		return err
	}
	category, err := h.services.Category.Create(c.Context(),input)
	if err != nil {
		return err
	}
	return c.JSON(category)
}

func (h * CategorynHandler) GetAllCategories(c *fiber.Ctx) error {
	var inputs dto.GetListCategories
	inputs.Limit = 10
	inputs.Offset = 0
	if err := c.QueryParser(&inputs); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	categories,totalCount,err := h.services.Category.GetAll(c.Context(), inputs.SearchString,inputs.Limit,inputs.Offset)

	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	categoriesCopy := make([]dto.Category, len(categories))
	for i, category := range categories{
		categoriesCopy[i] = dto.Category{ID: category.ID, Name: category.Name}
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(inputs.Limit)))
	response := dto.GetListCategoriesResponse{
		Items: categoriesCopy,
		Total: totalPages,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

