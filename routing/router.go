package routing

import (
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

	app.Get("/", urlshortener.DashboardHandler(urlShortenerRepository))
	app.Get("/healthcheck", HealthCheckHandler())
	app.Get("/protected", users.IsAuthenticatedHandler(store, auth), ProtectedHandler())

	app.Get("/user", users.UserHandler(store))
	app.Get("/user/register", users.RegisterHandler(store, auth))
	app.Get("/user/login", users.LoginHandler(store, auth))
	app.Get("/user/logout", users.LogoutHandler(store, auth))
	app.Get("/user/callback", users.CallbackHandler(store, auth, userRepository))

	urlShortener := urlshortener.NewUrlShortener(urlShortenerRepository)
	app.Get("/:shortUrl", urlshortener.RedirectShortUrlHandler(urlShortener))
	app.Delete("/shortUrl/:shortUrl", users.IsAuthenticatedHandler(store, auth), urlshortener.DeleteShortHandler(store))
	app.Post("/makeShort", users.IsAuthenticatedHandler(store, auth), urlshortener.MakeShortHandler(store, urlShortener))
}
