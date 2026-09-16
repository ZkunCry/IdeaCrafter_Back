package auth

import (
	"startup_back/internal/user"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}


func (h *Handler) SignUp(c *fiber.Ctx) error {
	var input user.CreateUserInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	result, err := h.service.SignUp(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	setAuthCookies(c, result.AccessToken, result.RefreshToken)
	return c.Status(fiber.StatusCreated).JSON(user.NewResponse(result.User))
}


func (h *Handler) SignIn(c *fiber.Ctx) error {
	var input user.CreateUserInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	result, err := h.service.SignIn(c.Context(), input.Email, input.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	setAuthCookies(c, result.AccessToken, result.RefreshToken)
	return c.Status(fiber.StatusOK).JSON(user.NewResponse(result.User))
}


func (h *Handler) IdentityMe(c *fiber.Ctx) error {
	result, err := h.service.IdentityMe(c.Context(), c.Cookies("access_token"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(user.NewResponse(result.User))
}


func (h *Handler) Refresh(c *fiber.Ctx) error {
	accessToken, err := h.service.RefreshToken(c.Context(), c.Cookies("refresh_token"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/",
		MaxAge:   900,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"access_token": accessToken})
}


func (h *Handler) LogOut(c *fiber.Ctx) error {
	c.ClearCookie("access_token")
	c.ClearCookie("refresh_token")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
}

func setAuthCookies(c *fiber.Ctx, accessToken, refreshToken string) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/",
		MaxAge:   900,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
	})
}
