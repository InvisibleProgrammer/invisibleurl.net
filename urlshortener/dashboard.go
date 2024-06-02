package urlshortener

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func DashboardHandler(urlShortenerRepository *UrlShortenerRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		allUrls, err := urlShortenerRepository.GetAll()
		if err != nil {
			c.SendStatus(fiber.StatusInternalServerError)
			return fmt.Errorf("couldn't receive dahboard items")
		}

		return c.Render("index", fiber.Map{
			"Title":     "InvisibleUrl.Net",
			"ShortURLs": allUrls,
		})
	}
}
