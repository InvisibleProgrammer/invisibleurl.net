package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func RegisterHandler(store *session.Store) fiber.Handler {

	return func(c *fiber.Ctx) error {

		// Todo: register user, send a confirmation email, redirect to a welcome_confirmpage link
		return c.SendStatus(fiber.StatusCreated)
	}

}
