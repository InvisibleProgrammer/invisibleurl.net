package routing

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	auditlog "invisibleprogrammer.com/invisibleurl/audit_log"
	"invisibleprogrammer.com/invisibleurl/urlshortener"
	"invisibleprogrammer.com/invisibleurl/users"
)

func RegisterRoutes(
	app *fiber.App,
	store *session.Store,
	userRepository *users.UserRepository,
	urlShortenerRepository *urlshortener.UrlShortenerRepository,
	auditLogService *auditlog.AuditLogService) {

	app.Static("/", "./public", fiber.Static{
		Compress:      false,
		ByteRange:     true,
		Browse:        false,
		CacheDuration: 24 * time.Hour,
		MaxAge:        24 * 60 * 60, // 24 hours in seconds
	})

	app.Get("/", urlshortener.DashboardHandler(store, userRepository, urlShortenerRepository))
	app.Post("/filter/", users.IsAuthenticatedHandler(store), urlshortener.FilterHandler(store, userRepository, urlShortenerRepository))
	app.Get("/healthcheck", HealthCheckHandler())
	app.Get("/protected", users.IsAuthenticatedHandler(store), ProtectedHandler())

	app.Get("/user", users.UserHandler(store))
	app.Get("/user/sign-up", users.GetSignUpHandler())
	app.Post("/user/sign-up", users.PostSignUpHandler(store, userRepository, auditLogService))
	app.Get("/user/sign-in", users.GetSignInHandler())
	app.Post("/user/sign-in", users.SignInHandler(store, userRepository, auditLogService))
	app.Get("/user/activate/:activationTicket", users.ActivationHandler(store, userRepository, auditLogService))
	app.Get("/user/sign-out", users.SignOutHandler(store, userRepository, auditLogService))

	urlShortener := urlshortener.NewUrlShortener(urlShortenerRepository)
	app.Get("/:shortUrl", urlshortener.RedirectShortUrlHandler(urlShortener))
	app.Delete("/shortUrl/:shortUrl", users.IsAuthenticatedHandler(store), urlshortener.DeleteShortHandler(store, userRepository, urlShortenerRepository, urlShortener, auditLogService))
	app.Post("/makeShort", users.IsAuthenticatedHandler(store), urlshortener.MakeShortHandler(store, userRepository, urlShortenerRepository, urlShortener, auditLogService))
	app.Post("/makeCustomShort", users.IsAuthenticatedHandler(store), urlshortener.MakeShortHandler(store, userRepository, urlShortenerRepository, urlShortener, auditLogService))
}
