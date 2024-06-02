package urlshortener

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type ShortenerHandler struct {
	urlshortener *UrlShortener
}

func NewShortenerHandler(urlshortener UrlShortener) *ShortenerHandler {
	return &ShortenerHandler{
		urlshortener: &urlshortener,
	}
}

func MakeShortHandler(store *session.Store, urlShortener *UrlShortener) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session, err := store.Get(c)
		if err != nil {
			log.Fatalf("Couldn't receive session: %v", err)
		}

		userId := session.Get("userId")
		fullUrl := c.FormValue("fullUrl")

		shortUrl, err := urlShortener.MakeShortUrl(userId.(string), fullUrl)
		if err != nil {
			log.Printf("Error on shortening: %v", err)
		} else {
			log.Printf("Shortened version: %s", shortUrl)
		}

		return c.Redirect("/")
	}
}

func DeleteShortHandler(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {

		session, err := store.Get(c)
		if err != nil {
			log.Fatalf("Couldn't receive session: %v", err)
		}

		userId := session.Get("userId")

		shortUrl := c.Params("shortUrl")

		err = DeleteShortUrl(userId.(string), shortUrl)
		if err != nil {
			errorMessage := fmt.Sprintf("Error on deleting %s: %v", shortUrl, err)
			log.Print(errorMessage)

			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": errorMessage,
			})
		}

		return c.SendStatus(fiber.StatusOK)

	}
}

func RedirectShortUrlHandler(urlShortener *UrlShortener) fiber.Handler {
	return func(c *fiber.Ctx) error {

		short := c.Params("shortUrl")

		fullUrl, err := urlShortener.GetFullUrl(short)
		if err != nil {
			return c.SendString("waaat")
		}

		return c.Redirect(fullUrl, http.StatusFound)
	}
}
