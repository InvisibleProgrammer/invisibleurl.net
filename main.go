package main

import (
	"encoding/gob"
	"os"

	"golang.org/x/exp/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"invisibleprogrammer.com/invisibleurl/authenticator"
	repository "invisibleprogrammer.com/invisibleurl/db"
	"invisibleprogrammer.com/invisibleurl/routing"
	"invisibleprogrammer.com/invisibleurl/urlshortener"
	"invisibleprogrammer.com/invisibleurl/users"
)

func main() {

	handler := slog.NewJSONHandler(os.Stdout, nil)
	log := slog.New(handler)

	log.Info("Starting Fiber application",
		slog.String("version", "v2.52.1"),
		slog.String("address", "127.0.0.1:3000"),
		slog.String("host", "0.0.0.0"),
		slog.Int("port", 3000),
		slog.Int("pid", os.Getpid()))

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Error("Failed to load the env vars: %v", err)
	}
	godotenv.Load(".env")

	// Initialize repositories
	repository, err := repository.NewRepository()
	if err != nil {
		log.Error("Failed to initialize db: %v\n", err)
	}

	userRepository := users.NewUserRepository(repository)
	urlShortenerRepository := urlshortener.NewUrlShortenerRepository(repository)

	// Initialize server session
	store := session.New()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	// OAuth2 authenticator
	auth, err := authenticator.New()
	if err != nil {
		log.Error("Failed to initialize the authenticator: %v\n", err)
	}

	// HTML templates
	engine := html.New("./views", ".html")

	// Initialize engine
	app := fiber.New(fiber.Config{
		Views:             engine,
		ViewsLayout:       "layouts/main",
		PassLocalsToViews: true,
	})

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: `{"time":"${time}","level":"INFO","msg":"Incoming request","request.time":"${time}","request.method":"${method}","request.host":"${host}","request.path":"${path}","request.query":"${query}","request.params":"${params}","request.route":"${route}","request.ip":"${ip}","request.x-forwarded-for":"${ips}","request.referer":"${referer}","request.length":${bytesReceived},"response.time":"${time}","response.latency":"${latency}","response.status":${status},"response.length":${bytesSent},"id":"${locals:requestid}"}` + "\n",
	}))

	app.Use(logger.New(logger.Config{
		Format: `{"time":"${time}","level":"INFO","msg":"Incoming request","request.time":"${time}","request.method":"${method}","request.host":"${host}","request.path":"${path}","request.query":"${query}","request.params":"${params}","request.route":"${route}","request.ip":"${ip}","request.x-forwarded-for":"${ips}","request.referer":"${referer}","request.length":${bytesReceived},"response.time":"${time}","response.latency":"${latency}","response.status":${status},"response.length":${bytesSent},"id":"${locals:requestid}"}` + "\n",
	}))

	// Show authenticated user name on header partial
	users.RegisterUsernameMiddleware(app, store)

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "http://localhost:3000",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	// Set up routing
	routing.RegisterRoutes(app, store, auth, userRepository, urlShortenerRepository)

	log.Error("Failed to start server", slog.String("error", app.Listen(":3000").Error()))
}
