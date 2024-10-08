package users

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func IsAuthenticatedHandler(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session, err := store.Get(c)
		if err != nil {
			log.Fatalf("Couldn't receive session: %v", err)
		}

		if session.Get("publicId") == nil {
			return c.Redirect("/", http.StatusSeeOther)
		} else {
			return c.Next()
		}
	}
}

func RegisterUsernameMiddleware(app *fiber.App, store *session.Store) {
	app.Use(func(c *fiber.Ctx) error {

		session, err := store.Get(c)
		if err != nil {
			log.Fatalf("Couldn't receive sesion: %v", err)
		}

		c.Locals("emailAddress", session.Get("emailAddress"))

		return c.Next()
	})
}
