package middleware

import (
	"startup_back/internal/service"

	"github.com/gofiber/fiber/v2"
)


func RequireAuth(tokenService service.TokenService) fiber.Handler{
	return func(c *fiber.Ctx) error{
		accessToken := c.Cookies("access_token")
		if accessToken ==""{
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missiung access token"})
		}
		userID,err:= tokenService.ValidateAccessToken(accessToken)
		if err != nil{
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}
		c.Locals("user_id", userID)
		return c.Next()
	}
}