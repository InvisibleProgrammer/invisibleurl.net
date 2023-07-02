package userhandler

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"invisibleprogrammer.com/invisibleurl/authenticator"
)

func IsAuthenticatedHandler(store *session.Store, auth *authenticator.Authenticator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session, err := store.Get(c)
		if err != nil {
			log.Fatalf("Couldn't receive session: %v", err)
		}

		if session.Get("profile") == nil {
			return c.Redirect("/", http.StatusSeeOther)
		} else {
			return c.Next()
		}
	}

}
