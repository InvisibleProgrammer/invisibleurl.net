package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func SignOutHandler(store *session.Store) fiber.Handler {

	return func(c *fiber.Ctx) error {
		//Todo: remove info from session, redirect to start page

		return c.SendStatus(fiber.StatusOK)
	}

}
