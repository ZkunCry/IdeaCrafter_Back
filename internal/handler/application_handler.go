package handler

import (
	"fmt"
	"startup_back/internal/service"

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

	formData := c.FormValue("test")
	fmt.Println(formData)

	return nil
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
	return nil
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
	return nil
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
	return nil
}

