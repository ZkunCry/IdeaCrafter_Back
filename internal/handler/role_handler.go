package handler

import (
	"startup_back/internal/dto"
	"startup_back/internal/service"

	"github.com/gofiber/fiber/v2"
)
type RoleHandler struct {
	services *service.Services
}

func NewRoleHandler(services *service.Services) *RoleHandler {
	return &RoleHandler{services: services}
}


func (r *RoleHandler) CreateRole(c *fiber.Ctx) error {
	var inputs dto.CreateRoleInput

	if err:= c.BodyParser(&inputs); err !=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	role,err:= r.services.Role.Create(c.Context(), &inputs)

	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(role)
}