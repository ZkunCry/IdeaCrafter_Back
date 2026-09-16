package startup

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
	s3      *s3.Client
	bucket  string
}

func NewHandler(service Service, s3Client *s3.Client, bucket string) *Handler {
	return &Handler{service: service, s3: s3Client, bucket: bucket}
}

func actorID(c *fiber.Ctx) (uint, bool) {
	userID, ok := c.Locals("user_id").(uint)
	return userID, ok
}

func respondError(c *fiber.Ctx, err error) error {
	if errors.Is(err, ErrForbidden) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
}

// formPayload holds the multipart fields shared by create and update.
type formPayload struct {
	Name             string
	ShortDescription string
	Description      string
	TargetAudience   string
	Problem          string
	Solution         string
	StageID          uint
	CategoryIDs      []uint
	LogoURL          string
}

func parseCategoryIDs(raw string) ([]uint, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}

	parts := strings.Split(trimmed, ",")
	ids := make([]uint, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseUint(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid category id: %s", part)
		}
		ids = append(ids, uint(id))
	}
	return ids, nil
}


func (h *Handler) uploadLogo(c *fiber.Ctx) (string, error) {
	file, _ := c.FormFile("files")
	if file == nil {
		return "", nil
	}

	f, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("cannot open file")
	}
	defer f.Close()

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("cannot read file")
	}

	objectName := "photos/" + uuid.New().String() + filepath.Ext(file.Filename)
	if _, err = h.s3.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(h.bucket),
		Key:         aws.String(objectName),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(file.Header.Get("Content-Type")),
	}); err != nil {
		return "", fmt.Errorf("cannot upload file to S3: %w", err)
	}

	return fmt.Sprintf("https://storage.yandexcloud.net/%s/%s", h.bucket, objectName), nil
}


func (h *Handler) parseForm(c *fiber.Ctx, requireStage bool) (*formPayload, error) {
	payload := &formPayload{
		Name:             strings.TrimSpace(c.FormValue("name")),
		ShortDescription: strings.TrimSpace(c.FormValue("short_description")),
		Description:      strings.TrimSpace(c.FormValue("description")),
		TargetAudience:   strings.TrimSpace(c.FormValue("target_audience")),
		Problem:          strings.TrimSpace(c.FormValue("problem")),
		Solution:         strings.TrimSpace(c.FormValue("solution")),
	}

	rawStage := strings.TrimSpace(c.FormValue("stage_id"))
	if rawStage == "" {
		if requireStage {
			return nil, fmt.Errorf("stage_id is required")
		}
	} else {
		stageID, err := strconv.ParseUint(rawStage, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid stage_id")
		}
		payload.StageID = uint(stageID)
	}

	categoryIDs, err := parseCategoryIDs(c.FormValue("category_ids"))
	if err != nil {
		return nil, err
	}
	payload.CategoryIDs = categoryIDs

	logoURL, err := h.uploadLogo(c)
	if err != nil {
		return nil, err
	}
	payload.LogoURL = logoURL

	return payload, nil
}

func (h *Handler) CreateStartup(c *fiber.Ctx) error {
	userID, ok := actorID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	payload, err := h.parseForm(c, true)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	startup, err := h.service.Create(c.Context(), CreateInput{
		CreatorID:        userID,
		Name:             payload.Name,
		ShortDescription: payload.ShortDescription,
		Description:      payload.Description,
		TargetAudience:   payload.TargetAudience,
		Problem:          payload.Problem,
		Solution:         payload.Solution,
		StageID:          payload.StageID,
		CategoryIDs:      payload.CategoryIDs,
		LogoFile:         payload.LogoURL,
	})
	if err != nil {
		return respondError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(NewResponse(startup))
}

func (h *Handler) UpdateStartup(c *fiber.Ctx) error {
	userID, ok := actorID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	payload, err := h.parseForm(c, false)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	startup, err := h.service.Update(c.Context(), uint(id), UpdateInput{
		Name:             payload.Name,
		ShortDescription: payload.ShortDescription,
		Description:      payload.Description,
		TargetAudience:   payload.TargetAudience,
		Problem:          payload.Problem,
		Solution:         payload.Solution,
		StageID:          payload.StageID,
		CategoryIDs:      payload.CategoryIDs,
		LogoFile:         payload.LogoURL,
	}, userID)
	if err != nil {
		return respondError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(NewResponse(startup))
}

func (h *Handler) GetListStartups(c *fiber.Ctx) error {
	var input ListInput
	if err := c.QueryParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	if input.Limit <= 0 {
		input.Limit = 12
	}
	startups, totalCount, err := h.service.GetAll(c.Context(), input.SearchString, input.CategorySlug, input.Limit, input.Offset)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(input.Limit)))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"items": NewResponses(startups), "total_count": totalPages})
}

func (h *Handler) GetStartupByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	startup, err := h.service.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(NewResponse(startup))
}

func (h *Handler) GetUserStartups(c *fiber.Ctx) error {
	userID, ok := actorID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	startups, err := h.service.GetUserStartups(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	items := make([]Response, 0, len(startups))
	for i := range startups {
		items = append(items, NewResponse(&startups[i]))
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"items": items})
}

func (h *Handler) AddCategories(c *fiber.Ctx) error {
	startupID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid startup id"})
	}
	var input AddCategoriesInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	startup, err := h.service.AddCategories(c.Context(), uint(startupID), input)
	if err != nil {
		return respondError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(NewResponse(startup))
}
