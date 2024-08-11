package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func ActivationHandler(store *session.Store, userRepository *UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		activationTicket := c.Params("activationTicket")

		if activationTicket == "" {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		var userId int64
		var err error
		if userId, err = userRepository.Get_UserId_by_ActivationTicket(activationTicket); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err := userRepository.Activate_User(userId); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.Redirect("/", fiber.StatusFound)
	}
}
