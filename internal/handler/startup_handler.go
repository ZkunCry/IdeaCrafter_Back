package handler

import (
	"startup_back/internal/dto"
	"startup_back/internal/service"

	"github.com/gofiber/fiber/v2"
)

type StartupHandler struct {
	services *service.Services
}

func NewStartupHandler(services *service.Services) *StartupHandler {
	return &StartupHandler{services: services}
}

func (s *StartupHandler) CreateStartup(c *fiber.Ctx) error {
	var inputs dto.CreateStartupInput
	if err := c.BodyParser(&inputs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	startup,err := s.services.Startup.Create(c.Context(), &inputs, inputs.CategoryIDs, inputs.VacancyIDs)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(startup)
}