package routing

import "github.com/gofiber/fiber/v2"

func HealthCheckHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.SendString("ok")
		c.SendStatus(fiber.StatusOK)
		return nil
	}
}
