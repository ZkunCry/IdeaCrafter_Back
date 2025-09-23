package handler

import (
	"startup_back/internal/service"

	"github.com/gofiber/fiber/v2"
)


type ApplicationHandler struct {
	services *service.Services
}

func NewApplicationHandler(services *service.Services) *ApplicationHandler {
	return &ApplicationHandler{services: services}
}


func (r *ApplicationHandler) CreateApplication(c *fiber.Ctx) error {
	return nil
}
func (r * ApplicationHandler) UpdateApplication(c *fiber.Ctx) error {
	return nil
}

func (r * ApplicationHandler) UpdateApplicationStatus(c *fiber.Ctx) error {
	return nil
}
func (r *ApplicationHandler) GetApplicationByID(c *fiber.Ctx) error {
	return nil
}

