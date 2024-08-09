package users

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func LoginHandler(store *session.Store) fiber.Handler {

	return func(c *fiber.Ctx) error {

		state, err := generateRandomState()
		if err != nil {
			c.SendString(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		session, err := store.Get(c)
		if err != nil {
			log.Fatalf("Couldn't receive sesion: %v", err)
		}

		session.Set("state", state)

		if err := session.Save(); err != nil {
			c.SendString(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusOK)
	}

}
