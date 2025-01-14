package users

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func SignOutHandler(store *session.Store) fiber.Handler {

	return func(c *fiber.Ctx) error {

		session, err := store.Get(c)
		if err != nil {
			log.Fatalf("Couldn't receive sesion: %v", err)
		}

		session.Delete("state")
		session.Delete("userId")
		session.Delete("emailAddress")
		session.Delete("publicId")

		if err := session.Save(); err != nil {
			c.SendString(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.Redirect("/", fiber.StatusFound)
	}

}
