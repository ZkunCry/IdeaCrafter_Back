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
	userID := c.Locals("user_id").(uint) 
	startup,err := s.services.Startup.Create(c.Context(), &dto.CreateStartupInput{Name: inputs.Name, Description: inputs.Description, CreatorId: userID}, inputs.CategoryIDs)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(startup)
}

func (s * StartupHandler) GetListStartups(c * fiber.Ctx) error{
	var inputs dto.GetStartupList
	if err := c.QueryParser(&inputs); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	startups,err := s.services.Startup.List(c.Context(),inputs.Limit,inputs.Offset)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(startups)
}