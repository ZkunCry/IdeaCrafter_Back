package handler

import (
	"startup_back/internal/dto"
	"startup_back/internal/service"

	"github.com/gofiber/fiber/v2"
)
type StageHandler struct {
	services *service.Services
}

func NewStageHandler(services *service.Services) *StageHandler {
	return &StageHandler{services: services}
}
// CreateStage godoc
// @Summary      Создание этапа
// @Description  Создает новый этап стартапа (например, MVP, Идея)
// @Tags         stages
// @Accept       json
// @Produce      json
// @Param        stage  body      dto.CreateStageInput  true  "Данные нового этапа"
// @Success      200    {object}  dto.StageResponse
// @Failure      400    {object}  map[string]string
// @Router       /stage [post]
func (s *StageHandler) CreateStage(c *fiber.Ctx) error {
	var inputs dto.CreateStageInput

	if err := c.BodyParser(&inputs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	stage, err := s.services.Stage.Create(c.Context(), &inputs)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(stage)
}
// CreateStage godoc
// @Summary      Получение списка этапов
// @Description  Ищет список этапов (например, MVP, Идея)
// @Tags         stages
// @Accept       json
// @Produce      json
// @Success      200    {object}  []dto.StageResponse
// @Failure      400    {object}  map[string]string
// @Router       /stage [get]
func (s * StageHandler) GetList(c * fiber.Ctx) error{
	stages, err := s.services.Stage.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(stages)
}