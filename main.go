package main

import (
	"encoding/gob"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	"github.com/lmittmann/tint"
	slogfiber "github.com/samber/slog-fiber"
	auditlog "invisibleprogrammer.com/invisibleurl/audit_log"
	repository "invisibleprogrammer.com/invisibleurl/db"
	"invisibleprogrammer.com/invisibleurl/environment"
	"invisibleprogrammer.com/invisibleurl/routing"
	"invisibleprogrammer.com/invisibleurl/security"
	"invisibleprogrammer.com/invisibleurl/urlshortener"
	"invisibleprogrammer.com/invisibleurl/users"
)

func main() {

	logger := slog.New(tint.NewHandler(os.Stdout, nil))
	log := slog.New(logger.Handler())

	log.Info("Starting Fiber application",
		slog.String("version", "v2.52.1"),
		slog.String("address", "127.0.0.1:3000"),
		slog.String("host", "0.0.0.0"),
		slog.Int("port", 3000),
		slog.Int("pid", os.Getpid()))

	environment.Init()

	// Initialize repositories
	repository, err := repository.NewRepository()
	if err != nil {
		log.Error("Failed to initialize db: %v", slog.Any("err", err))
	}

	userRepository := users.NewUserRepository(repository)
	urlShortenerRepository := urlshortener.NewUrlShortenerRepository(repository)
	auditLogRepository := auditlog.NewAuditLogRepository(repository)

	// initialize audit logger
	auditLogger := auditlog.NewAuditLogService(auditLogRepository)

	// Initialize server session
	store := session.New()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	// HTML templates
	engine := html.New("./views", ".html")

	// Initialize engine
	app := fiber.New(fiber.Config{
		Views:             engine,
		ViewsLayout:       "layouts/main",
		PassLocalsToViews: true,
	})

	app.Use(requestid.New())

	app.Use(slogfiber.New(logger))

	// security
	security.RegisterRateLimitingMiddleware(app)

	// Show authenticated user name on header partial
	users.RegisterUsernameMiddleware(app, store)

	// Set up routing
	routing.RegisterRoutes(app, store, userRepository, urlShortenerRepository, auditLogger)

	log.Error("Failed to start server", slog.String("error", app.Listen("127.0.0.1:8080").Error()))
}
