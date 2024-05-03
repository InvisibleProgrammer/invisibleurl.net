package routing

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"invisibleprogrammer.com/invisibleurl/authenticator"
	"invisibleprogrammer.com/invisibleurl/shortenerhandler"
	"invisibleprogrammer.com/invisibleurl/urlshortener"
	"invisibleprogrammer.com/invisibleurl/userhandler"
)

func RegisterRoutes(app *fiber.App, store *session.Store, auth *authenticator.Authenticator) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":     "InvisibleUrl.Net",
			"ShortURLs": urlshortener.GetAll(),
		})
	})

	app.Get("/protected", userhandler.IsAuthenticatedHandler(store, auth), func(c *fiber.Ctx) error {
		return c.Render("protected", fiber.Map{
			"Title": "Hello, World. You shouldn't see that unless you logged in!",
		}, "layouts/main")
	})

	app.Get("/user", userhandler.UserHandler(store))
	app.Get("/user/login", userhandler.LoginHandler(store, auth))
	app.Get("/user/logout", userhandler.LogoutHandler(store, auth))
	app.Get("/user/callback", userhandler.CallbackHandler(store, auth))

	app.Get("/:shortUrl", func(c *fiber.Ctx) error {
		short := c.Params("shortUrl")

		fullUrl, err := urlshortener.GetFullUrl(short)
		if err != nil {
			return c.SendString("waaat")
		}

		return c.Redirect(fullUrl, http.StatusFound)
	})

	app.Delete("/shortUrl/:shortUrl", userhandler.IsAuthenticatedHandler(store, auth), shortenerhandler.DeleteShortHandler(store))

	app.Post("/makeShort", userhandler.IsAuthenticatedHandler(store, auth), shortenerhandler.MakeShortHandler(store))

}
