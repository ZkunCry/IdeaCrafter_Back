package handler

import (
	"startup_back/internal/dto"
	"startup_back/internal/service"

	"github.com/gofiber/fiber/v2"
)

type VacancyHandler struct {
	services *service.Services
}
func NewVacancyHandler(services *service.Services) *VacancyHandler {
	return &VacancyHandler{services: services}
}


func (h * VacancyHandler) CreateVacancy(c *fiber.Ctx) error {
	var inputs dto.CreateVacancyInput

	if err:= c.BodyParser(&inputs); err !=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	vacancy,err:= h.services.Vacancy.Create(c.Context(), &inputs)

	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(vacancy)
}



func (h *VacancyHandler) GetVacancyByID(c *fiber.Ctx) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	vacancy, err := h.services.Vacancy.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	vacancyResponse := dto.VacancyResponse{
		ID:          vacancy.ID,
		RoleName:    vacancy.Role.Name,
		Description: vacancy.Description,
		StartupID:   vacancy.StartupID,
		IsOpen:      vacancy.IsOpen,
		Role:        vacancy.Role,
	}
	return c.JSON(vacancyResponse)
}


func (h *VacancyHandler) GetVacanciesByStartup(c *fiber.Ctx) error {
	startupID, err := parseUintParam(c, "startupID")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid startup id"})
	}

	vacancies, err := h.services.Vacancy.GetByStartupID(c.Context(), startupID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(vacancies)
}

func (h *VacancyHandler) UpdateVacancy(c *fiber.Ctx) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var input dto.UpdateVacancyInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	vacancy,err := h.services.Vacancy.Update(c.Context(), id, &input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(vacancy)
}