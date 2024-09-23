package urlshortener

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"invisibleprogrammer.com/invisibleurl/users"
)

type ShortenerHandler struct {
	urlshortener *UrlShortener
}

func NewShortenerHandler(urlshortener UrlShortener) *ShortenerHandler {
	return &ShortenerHandler{
		urlshortener: &urlshortener,
	}
}

func MakeShortHandler(store *session.Store, userRepository *users.UserRepository, urlShortenerRepostiory *UrlShortenerRepository, urlShortener *UrlShortener) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session, err := store.Get(c)
		if err != nil {
			log.Fatalf("Couldn't receive session: %v", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		publicId := session.Get("publicId").(string)
		user, err := userRepository.Get_UserId_by_PublicId(publicId)
		if err != nil {
			log.Printf("Cannot get user by public id")
			return c.SendStatus(fiber.StatusBadRequest)
		}

		fullUrl := c.FormValue("fullUrl")

		shortUrlId, err := urlShortenerRepostiory.GetNextUrlId()
		if err != nil {
			log.Printf("Error on shortening: %v", err)
			return c.SendStatus(fiber.StatusBadRequest)

		}

		shortUrl, err := urlShortener.MakeShortUrl(shortUrlId)
		if err != nil {
			log.Printf("Error on shortening: %v", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		shortenedUrl := ShortenedUrl{
			UrlId:    int(shortUrlId),
			UserId:   user.Id,
			ShortUrl: shortUrl,
			FullUrl:  fullUrl,
		}

		err = urlShortenerRepostiory.Store(shortenedUrl)
		if err != nil {
			log.Fatalf("error on storing shortened url. UserId: %d, Full url: %s, short url: %s, error: %v", user.Id, fullUrl, shortUrl, err)
			return c.SendStatus(fiber.StatusBadRequest)
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
