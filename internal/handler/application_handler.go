package handler

import (
	"startup_back/internal/dto"
	"startup_back/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)


type ApplicationHandler struct {
	services *service.Services
}

func NewApplicationHandler(services *service.Services) *ApplicationHandler {
	return &ApplicationHandler{services: services}
}


// CreateApplication
// @Summary Create a new application
// @Description Создает новую заявку на вакансию
// @Tags applications
// @Accept json
// @Produce json
// @Param input body dto.CreateApplicationInput true "Application info"
// @Success 201 {object} dto.CreateApplicationInput
// @Failure 400 {object} map[string]string
// @Router /application [post]
func (r *ApplicationHandler) CreateApplication(c *fiber.Ctx) error {

	var input dto.CreateApplicationInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	input.UserID = userID

	application, err := r.services.Application.Create(c.Context(), &input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(application)
}
// UpdateApplication
// @Summary Update an application
// @Description Обновляет данные заявки (например, сообщение)
// @Tags applications
// @Accept json
// @Produce json
// @Param id path int true "Application ID"
// @Param input body dto.CreateApplicationInput true "Updated application info"
// @Success 200 {object} dto.CreateApplicationInput
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /application/{id} [put]
func (r * ApplicationHandler) UpdateApplication(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var input dto.UpdateApplicationInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	application, err := r.services.Application.Update(c.Context(), uint(id), &input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(application)
}
// UpdateApplicationStatus
// @Summary Update application status
// @Description Изменяет статус заявки: pending, accepted, rejected
// @Tags applications
// @Accept json
// @Produce json
// @Param id path int true "Application ID"
// @Param input body dto.UpdateApplicationStatusInput true "New status"
// @Success 200 {object} dto.UpdateApplicationStatusInput
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /applications/{id}/status [patch]
func (r * ApplicationHandler) UpdateApplicationStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var input dto.UpdateApplicationStatusInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	application, err := r.services.Application.UpdateStatus(c.Context(), uint(id), &input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(application)
}

// GetApplicationByID
// @Summary Get application by ID
// @Description Получает заявку по ID
// @Tags applications
// @Accept json
// @Produce json
// @Param id path int true "Application ID"
// @Success 200 {object} dto.CreateApplicationInput
// @Failure 404 {object} map[string]string
// @Router /application/{id} [get]
func (r *ApplicationHandler) GetApplicationByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	application, err := r.services.Application.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(application)
}

