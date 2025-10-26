package handler

import (
	"math"
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
	var input dto.CreateStartupInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid input format",
		})
	}


	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}


	input.CreatorID = userID


	startup, err := s.services.Startup.Create(c.Context(), &input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}


	response := dto.StartupResponse{
		ID:          startup.ID,
		Name:        startup.Name,
		Description: startup.Description,
		TargetAudience: startup.TargetAudience,
		Solution: startup.Solution,
		ShortDescription: startup.ShortDescription,
		Creator:     startup.Creator,
		Problem: startup.Problem,
		Categories:  startup.Categories,
		Files:       startup.Files,
		Vacansies:   startup.Vacancies,
		Stage: startup.Stage,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}


func (s * StartupHandler) GetListStartups(c * fiber.Ctx) error{
	var inputs dto.GetStartupList
	if err := c.QueryParser(&inputs); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	startups,totalCount,err := s.services.Startup.GetAll(c.Context(),inputs.Limit,inputs.Offset)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	startupResponse:=[]dto.StartupResponse{}
	for _,startup := range startups{
		startupResponse = append(startupResponse, dto.StartupResponse{
			ID:          startup.ID,
		Name:        startup.Name,
		Description: startup.Description,
		TargetAudience: startup.TargetAudience,
		Solution: startup.Solution,
		ShortDescription: startup.ShortDescription,
		Creator:     startup.Creator,
		Problem: startup.Problem,
		Categories:  startup.Categories,
		Files:       startup.Files,
		Vacansies:   startup.Vacancies,
		Stage: startup.Stage,
		})
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(inputs.Limit)))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"startups": startupResponse,
		"total_count": totalPages,
	})
}