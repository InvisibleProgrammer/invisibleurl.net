package routing

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"invisibleprogrammer.com/invisibleurl/urlshortener"
	"invisibleprogrammer.com/invisibleurl/users"
)

func RegisterRoutes(
	app *fiber.App,
	store *session.Store,
	userRepository *users.UserRepository,
	urlShortenerRepository *urlshortener.UrlShortenerRepository) {

	app.Get("/", urlshortener.DashboardHandler(urlShortenerRepository))
	app.Get("/healthcheck", HealthCheckHandler())
	app.Get("/protected", users.IsAuthenticatedHandler(store), ProtectedHandler())

	app.Get("/user", users.UserHandler(store))
	app.Get("/user/sign-up", users.GetSignUpHandler())
	app.Post("/user/sign-up", users.PostSignUpHHandler(store, userRepository))
	app.Get("/user/login", users.LoginHandler(store))
	app.Get("/user/logout", users.LogoutHandler(store))

	urlShortener := urlshortener.NewUrlShortener(urlShortenerRepository)
	app.Get("/:shortUrl", urlshortener.RedirectShortUrlHandler(urlShortener))
	app.Delete("/shortUrl/:shortUrl", users.IsAuthenticatedHandler(store), urlshortener.DeleteShortHandler(store))
	app.Post("/makeShort", users.IsAuthenticatedHandler(store), urlshortener.MakeShortHandler(store, urlShortener))
}
