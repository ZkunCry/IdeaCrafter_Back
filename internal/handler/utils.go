package handler

import "github.com/gofiber/fiber/v2"
func parseUintParam(c *fiber.Ctx, name string) (uint, error) {
	id, err := c.ParamsInt(name)
	if err != nil || id <= 0 {
		return 0, fiber.NewError(fiber.StatusBadRequest, "invalid "+name)
	}
	return uint(id), nil
}