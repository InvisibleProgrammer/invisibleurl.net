package routing

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"invisibleprogrammer.com/invisibleurl/authenticator"
	"invisibleprogrammer.com/invisibleurl/urlshortener"
	"invisibleprogrammer.com/invisibleurl/users"
)

func RegisterRoutes(
	app *fiber.App,
	store *session.Store,
	auth *authenticator.Authenticator,
	userRepository *users.UserRepository,
	urlShortenerRepository *urlshortener.UrlShortenerRepository) {

	app.Get("/", func(c *fiber.Ctx) error {
		allUrls, err := urlShortenerRepository.GetAll()
		if err != nil {
			c.SendStatus(fiber.StatusInternalServerError)
			return fmt.Errorf("couldn't receive dahboard items")
		}

		return c.Render("index", fiber.Map{
			"Title":     "InvisibleUrl.Net",
			"ShortURLs": allUrls,
		})
	})

	app.Get("/protected", users.IsAuthenticatedHandler(store, auth), func(c *fiber.Ctx) error {
		return c.Render("protected", fiber.Map{
			"Title": "Hello, World. You shouldn't see that unless you logged in!",
		}, "layouts/main")
	})

	app.Get("/user", users.UserHandler(store))
	app.Get("/user/login", users.LoginHandler(store, auth))
	app.Get("/user/logout", users.LogoutHandler(store, auth))
	app.Get("/user/callback", users.CallbackHandler(store, auth, userRepository))

	urlShortener := urlshortener.NewUrlShortener(urlShortenerRepository)
	app.Get("/:shortUrl", urlshortener.RedirectShortUrlHandler(urlShortener))
	app.Delete("/shortUrl/:shortUrl", users.IsAuthenticatedHandler(store, auth), urlshortener.DeleteShortHandler(store))
	app.Post("/makeShort", users.IsAuthenticatedHandler(store, auth), urlshortener.MakeShortHandler(store, urlShortener))

}
