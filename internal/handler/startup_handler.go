package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"startup_back/internal/dto"
	"startup_back/internal/service"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StartupHandler struct {
	services *service.Services
}

func NewStartupHandler(services *service.Services) *StartupHandler {
	return &StartupHandler{services: services}
}
// CreateStartup godoc
// @Summary      Создать стартап
// @Description  Создает новый стартап от имени авторизованного пользователя
// @Tags         startups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        startup  body      dto.CreateStartupInput  true  "Данные стартапа"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Router       /startups [post]
func (s *StartupHandler) CreateStartup(c *fiber.Ctx) error {
	var input dto.CreateStartupInput

	input.Name = c.FormValue("name")
  input.ShortDescription = c.FormValue("short_description")
  input.Description = c.FormValue("description")
  input.TargetAudience = c.FormValue("target_audience")
  input.Problem = c.FormValue("problem")
  input.Solution = c.FormValue("solution")
	stageID,err := strconv.ParseInt(c.FormValue("stage_id"),10,64) 
	fmt.Println(stageID)
	input.StageID = uint(stageID)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid input stage",
			"type" : err.Error(),	
		})
	}
	categoryIDs := strings.Split(c.FormValue("category_ids"), ",")
	if len(categoryIDs) >1 {
		for _, categoryID := range categoryIDs {
			id, err := strconv.ParseUint(categoryID, 10, 64)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "invalid input",
				})
			}
			input.CategoryIDs = append(input.CategoryIDs, uint(id))
		}
	}
	fmt.Println(len(input.CategoryIDs))
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}
	input.CreatorID = userID
	file,_ := c.FormFile("files")

	if file != nil{
		f,err := file.Open()
		if err != nil{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot open file",
			})

		}
		defer f.Close()
		fileBytes, err := io.ReadAll(f)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "cannot read file")
    }
		objectName := "photos/" + uuid.New().String() + filepath.Ext(file.Filename)
		_, err = s.services.S3.PutObject(
        context.Background(),
        &s3.PutObjectInput{
            Bucket:      aws.String("idea-crafter"),
            Key:         aws.String(objectName),
            Body:        bytes.NewReader(fileBytes),
            ContentType: aws.String(file.Header.Get("Content-Type")),
        },
    )
		if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "cannot upload file to S3: "+err.Error())
    }

		input.LogoFile = fmt.Sprintf(
        "https://storage.yandexcloud.net/%s/%s",
        "idea-crafter",
        objectName,
    )
		fmt.Printf("LOGO URL: %s",input.LogoFile)
	}
	
	startup, err := s.services.Startup.Create(c.Context(), input)
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
		Creator:     dto.UserResponse{ID: startup.CreatorID, Username: startup.Creator.Username, Email: startup.Creator.Email}, 
		Problem: startup.Problem,
		Categories:  startup.Categories,
		Files:       startup.Files,
		Vacansies:   startup.Vacancies,
		Stage: dto.StageResponse{
			ID:   startup.StageID,
			Name: startup.Stage.Name,
		},
		LogoUrl: startup.LogoURL,
		
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// CreateStartup godoc
// @Summary      Получить стартап
// @Description  Получение стартапа по offset и limit
// @Tags         startups
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        startup  body      dto.CreateStartupInput  true  "Данные стартапа"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Router       /startups [get]
func (s * StartupHandler) GetListStartups(c * fiber.Ctx) error{
	var inputs dto.GetStartupList
	if err := c.QueryParser(&inputs); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	fmt.Println(inputs.SearchString)
	startups,totalCount,err := s.services.Startup.GetAll(c.Context(), inputs.SearchString,inputs.Limit,inputs.Offset)
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
		Creator:     dto.UserResponse{ID: startup.CreatorID, Username: startup.Creator.Username, Email: startup.Creator.Email},
		Problem: startup.Problem,
		Categories:  startup.Categories,
		Files:       startup.Files,
		Vacansies:   startup.Vacancies,
		Stage: dto.StageResponse{
			ID:   startup.StageID,
			Name: startup.Stage.Name,
		},
		LogoUrl: startup.LogoURL,
		})
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(inputs.Limit)))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"items": startupResponse,
		"total_count": totalPages,
	})
}

func (s * StartupHandler) GetStartupByID(c * fiber.Ctx) error{
	id, err := strconv.ParseUint(c.Params("id"),10,64)
	if err != nil{
		
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	startup,err := s.services.Startup.GetByID(c.Context(), uint(id))
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(startup)
}

func (s * StartupHandler) GetUserStartups(c * fiber.Ctx) error{
	fmt.Println("We are here")
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}
	startups,err := s.services.Startup.GetUserStartups(c.Context(),userID)
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
		Creator:     dto.UserResponse{ID: startup.CreatorID, Username: startup.Creator.Username, Email: startup.Creator.Email},
		Problem: startup.Problem,
		Categories:  startup.Categories,
		Files:       startup.Files,
		Vacansies:   startup.Vacancies,
		Stage: dto.StageResponse{
			ID:   startup.StageID,
			Name: startup.Stage.Name,
		},
		LogoUrl: startup.LogoURL,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"items": startupResponse,
	})
}