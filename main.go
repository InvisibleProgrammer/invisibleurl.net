package main

import (
	"encoding/gob"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
	slogfiber "github.com/samber/slog-fiber"
	repository "invisibleprogrammer.com/invisibleurl/db"
	"invisibleprogrammer.com/invisibleurl/routing"
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

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Error("Failed to load the env vars: %v", slog.Any("err", err))
	}
	godotenv.Load(".env")

	// Initialize repositories
	repository, err := repository.NewRepository()
	if err != nil {
		log.Error("Failed to initialize db: %v", slog.Any("err", err))
	}

	userRepository := users.NewUserRepository(repository)
	urlShortenerRepository := urlshortener.NewUrlShortenerRepository(repository)

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

	// Show authenticated user name on header partial
	users.RegisterUsernameMiddleware(app, store)

	// app.Use(cors.New(cors.Config{
	// 	AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
	// 	AllowOrigins:     "http://localhost:3000",
	// 	AllowCredentials: true,
	// 	AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	// }))

	// Set up routing
	routing.RegisterRoutes(app, store, userRepository, urlShortenerRepository)

	log.Error("Failed to start server", slog.String("error", app.Listen(":3000").Error()))
}
