package favorite

import (
	"errors"
	"strconv"

	"startup_back/internal/startup"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func actorID(c *fiber.Ctx) (uint, bool) {
	userID, ok := c.Locals("user_id").(uint)
	return userID, ok
}

func startupIDParam(c *fiber.Ctx) (uint, error) {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	return uint(id), err
}

func (h *Handler) AddFavorite(c *fiber.Ctx) error {
	userID, ok := actorID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	startupID, err := startupIDParam(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid startup id"})
	}
	count, err := h.service.Add(c.Context(), userID, startupID)
	if err != nil {
		if errors.Is(err, ErrStartupNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"startup_id": startupID, "is_favorite": true, "count": count})
}

func (h *Handler) RemoveFavorite(c *fiber.Ctx) error {
	userID, ok := actorID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	startupID, err := startupIDParam(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid startup id"})
	}
	count, err := h.service.Remove(c.Context(), userID, startupID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"startup_id": startupID, "is_favorite": false, "count": count})
}

func (h *Handler) GetMyFavoriteIDs(c *fiber.Ctx) error {
	userID, ok := actorID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	ids, err := h.service.GetStartupIDs(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"items": ids})
}

func (h *Handler) GetMyFavorites(c *fiber.Ctx) error {
	userID, ok := actorID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	startups, err := h.service.GetStartups(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"items": startup.NewResponses(startups)})
}

func (h *Handler) GetStartupFavoritesCount(c *fiber.Ctx) error {
	startupID, err := startupIDParam(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid startup id"})
	}
	count, err := h.service.CountByStartup(c.Context(), startupID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"startup_id": startupID, "count": count})
}
