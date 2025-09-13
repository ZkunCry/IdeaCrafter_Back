package handler

import (
	"startup_back/internal/service"

	"github.com/gofiber/fiber/v2"
)
type UserHandler struct{
	service service.UserService
}

func NewUserHandler(service service.UserService) * UserHandler{
	return &UserHandler{service:service}
}
func (h *UserHandler) Register(c *fiber.Ctx) error {
  var input service.CreateUserInput
  if err := c.BodyParser(&input); err != nil {
      return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
  }
  user, err := h.service.CreateUser(c.Context(), input)
  if err != nil {
      return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
  }

  return c.Status(fiber.StatusCreated).JSON(user)
}