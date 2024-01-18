package shortenerhandler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"invisibleprogrammer.com/invisibleurl/urlshortener"
)

func MakeShortHandler(store *session.Store) fiber.Handler {
	log.Printf("MakeShortHandler begin: \n")

	return func(c *fiber.Ctx) error {
		session, err := store.Get(c)
		if err != nil {
			log.Fatalf("Couldn't receive session: %v", err)
		}

		userId := session.Get("userId")
		fullUrl := c.FormValue("fullUrl")

		shortUrl, err := urlshortener.MakeShortUrl(userId.(string), fullUrl)
		if err != nil {
			log.Printf("Error on shortening: %v", err)
		} else {
			log.Printf("Shortened version: %s", shortUrl)
		}

		return c.Redirect("/")
	}
}
