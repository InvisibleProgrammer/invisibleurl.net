package urlshortener

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	auditlog "invisibleprogrammer.com/invisibleurl/audit_log"
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

func FilterHandler(store *session.Store, userRepository *users.UserRepository, urlShortenerRepository *UrlShortenerRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session, err := store.Get(c)
		if err != nil {
			log.Fatalf("Couldn't receive session: %v", err)
		}

		searchString := c.FormValue("search")

		publicId := session.Get("publicId")
		log.Printf("PublicId: %v\n", publicId)

		user, err := userRepository.Get_UserId_by_PublicId(publicId.(string))
		if err != nil {
			errorMessage := fmt.Sprintf("Cannot get user by public id: %s", err)
			log.Print(errorMessage)

			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": errorMessage,
			})
		}

		allUrls, err := urlShortenerRepository.GetPage(user.Id, 1, searchString)
		log.Printf("allUrls count: %d\n", len(allUrls))

		if err != nil {
			c.SendStatus(fiber.StatusInternalServerError)
			return fmt.Errorf("couldn't receive dahboard items: %v", err)
		}

		return c.Render("partials/manageurls", fiber.Map{
			"Title":     "InvisibleUrl.Net",
			"ShortURLs": allUrls,
		}, "layouts/empty")
	}
}

func DashboardHandler(store *session.Store, userRepository *users.UserRepository, urlShortenerRepository *UrlShortenerRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session, err := store.Get(c)
		if err != nil {
			log.Fatalf("Couldn't receive session: %v", err)
		}

		searchString := c.FormValue("search")

		publicId := session.Get("publicId")
		log.Printf("PublicId: %v\n", publicId)

		if publicId != nil {
			user, err := userRepository.Get_UserId_by_PublicId(publicId.(string))
			if err != nil {
				errorMessage := fmt.Sprintf("Cannot get user by public id: %s", err)
				log.Print(errorMessage)

				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": errorMessage,
				})
			}

			allUrls, err := urlShortenerRepository.GetPage(user.Id, 1, searchString)
			log.Printf("allUrls count: %d\n", len(allUrls))

			if err != nil {
				c.SendStatus(fiber.StatusInternalServerError)
				return fmt.Errorf("couldn't receive dahboard items: %v", err)
			}

			return c.Render("index", fiber.Map{
				"Title":     "InvisibleUrl.Net",
				"ShortURLs": allUrls,
			})
		}

		allUrls, err := urlShortenerRepository.GetDashboard()
		if err != nil {
			c.SendStatus(fiber.StatusInternalServerError)
			return fmt.Errorf("couldn't receive dahboard items")
		}

		return c.Render("index", fiber.Map{
			"ShortURLs": allUrls,
		})
	}
}

func MakeShortHandler(store *session.Store, userRepository *users.UserRepository, urlShortenerRepostiory *UrlShortenerRepository, urlShortener *UrlShortener, auditLogService *auditlog.AuditLogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auditLogEvent := auditlog.CREATE_CUSTOM_SHORT_URL

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
		shortUrl := c.FormValue("shortUrl")
		if len(shortUrl) > 50 {
			log.Printf("Short url is too long: %s", shortUrl)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		shortUrlId, err := urlShortenerRepostiory.GetNextUrlId()
		if err != nil {
			log.Printf("Error on shortening: %v", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if len(shortUrl) == 0 {
			shortUrl, err = urlShortener.MakeShortUrl(shortUrlId)

			if err != nil {
				log.Printf("Error on shortening: %v", err)
				return c.SendStatus(fiber.StatusBadRequest)
			}

			auditLogEvent = int(auditlog.CREATE_SHORT_URL)
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

		auditLogService.LogEvent(auditlog.Action(auditLogEvent), user.Id, c.Context().RemoteIP())
		return c.Redirect("/")
	}
}

func DeleteShortHandler(store *session.Store, userRepository *users.UserRepository, urlShoUrlShortenerRepository *UrlShortenerRepository, urlShortener *UrlShortener, auditLogService *auditlog.AuditLogService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		session, err := store.Get(c)
		if err != nil {
			log.Fatalf("Couldn't receive session: %v", err)
		}

		shortUrl := c.Params("shortUrl")
		publicId := session.Get("publicId")

		user, err := userRepository.Get_UserId_by_PublicId(publicId.(string))
		if err != nil {
			errorMessage := fmt.Sprintf("Error on deleting %s: %v", shortUrl, err)
			log.Print(errorMessage)

			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": errorMessage,
			})
		}

		err = urlShoUrlShortenerRepository.DeleteShortUrl(user.Id, shortUrl)
		if err != nil {
			errorMessage := fmt.Sprintf("Error on deleting %s: %v", shortUrl, err)
			log.Print(errorMessage)

			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": errorMessage,
			})
		}

		auditLogService.LogEvent(auditlog.DELETE_SHORTENED_URL, user.Id, c.Context().RemoteIP())
		c.SendStatus(fiber.StatusOK)
		return c.SendString("")
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
