package handler

import (
	"fmt"
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
// CreateVacancy godoc
// @Summary      Создать вакансию
// @Description  Создает новую вакансию для стартапа
// @Tags         vacancies
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        vacancy  body      dto.CreateVacancyInput  true  "Данные вакансии"
// @Success      200      {object}  dto.VacancyResponse
// @Failure      400      {object}  map[string]string
// @Router       /vacancy [post]
func (h * VacancyHandler) CreateVacancy(c *fiber.Ctx) error {
	var inputs dto.CreateVacancyInput

	if err:= c.BodyParser(&inputs); err !=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	vacancy,err:= h.services.Vacancy.Create(c.Context(), &inputs)

	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	fmt.Println(vacancy)
	role := dto.RoleResponse{
		ID:   vacancy.Role.ID,
		Name: vacancy.Role.Name,
	}
	vacancyResponse :=	dto.VacancyResponse{
		ID:          vacancy.ID,
		RoleName:    vacancy.Role.Name,
		Description: vacancy.Description,
		StartupID:   vacancy.StartupID,
		IsOpen:      vacancy.IsOpen,
		Role:        role,
	}

	return c.Status(fiber.StatusOK).JSON(vacancyResponse)
}

// GetVacancyByID godoc
// @Summary      Получить вакансию по ID
// @Description  Возвращает вакансию по её идентификатору
// @Tags         vacancies
// @Produce      json
// @Param        id   path      int  true  "ID вакансии"
// @Success      200  {object}  dto.VacancyResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router      /vacancy/{id} [get]
func (h *VacancyHandler) GetVacancyByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	vacancy, err := h.services.Vacancy.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	role :=dto.RoleResponse{
		ID:   vacancy.Role.ID,
		Name: vacancy.Role.Name,
	}
	vacancyResponse := dto.VacancyResponse{
		ID:          vacancy.ID,
		RoleName:    vacancy.Role.Name,
		Description: vacancy.Description,
		StartupID:   vacancy.StartupID,
		IsOpen:      vacancy.IsOpen,
		Role:        role,
	}
	return c.JSON(vacancyResponse)
}

// GetVacanciesByStartup godoc
// @Summary      Получить вакансии стартапа
// @Description  Возвращает список вакансий, связанных с конкретным стартапом
// @Tags         vacancies
// @Produce      json
// @Param        startupID  path  int  true  "ID стартапа"
// @Success      200        {array}   dto.VacancyResponse
// @Failure      400        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Router       /vacancy/startup/:startupID [get]
func (h *VacancyHandler) GetVacanciesByStartup(c *fiber.Ctx) error {
	startupID, err :=c.ParamsInt("id")
	fmt.Print(startupID)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid startup id"})
	}

	vacancies, err := h.services.Vacancy.GetByStartupID(c.Context(), uint(startupID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(vacancies)
}
// UpdateVacancy godoc
// @Summary      Обновить вакансию
// @Description  Обновляет данные существующей вакансии по ID
// @Tags         vacancies
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id       path      int                     true  "ID вакансии"
// @Param        vacancy  body      dto.UpdateVacancyInput  true  "Данные для обновления"
// @Success      200      {object}  dto.VacancyResponse
// @Failure      400      {object}  map[string]string
// @Router       /vacancy/{id} [put]
func (h *VacancyHandler) UpdateVacancy(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var input dto.UpdateVacancyInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	vacancy,err := h.services.Vacancy.Update(c.Context(), uint(id), &input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(vacancy)
}
// DeleteVacancy godoc
// @Summary      Удалить вакансию
// @Description  Удаляет вакансию по заданному id
// @Tags         vacancies
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id       path      int                     true  "ID вакансии"
// @Success      200  {object}  map[string]string      "Пример ответа"  example({"message":"vacancy deleted"})
// @Failure      400      {object}  map[string]string
// @Router       /vacancy/{id} [delete]
func (h *VacancyHandler) DeleteVacancy(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	err = h.services.Vacancy.Delete(c.Context(),uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "vacancy deleted"})
}