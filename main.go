package main

import (
	"encoding/gob"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Panicf("Failed to load the env vars: %v", err)
	}
	godotenv.Load(".env")

	// Initialize repositories
	repository, err := repository.NewRepository()
	if err != nil {
		log.Fatalf("Failed to initialize db: %v\n", err)
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
		log.Fatalf("Failed to initialize the authenticator: %v\n", err)
	}

	// HTML templates
	engine := html.New("./views", ".html")

	// Initialize engine
	app := fiber.New(fiber.Config{
		Views:             engine,
		ViewsLayout:       "layouts/main",
		PassLocalsToViews: true,
	})

	// Show authenticated user name on header partial
	users.RegisterUsernameMiddleware(app, store)

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	// Set up routing
	routing.RegisterRoutes(app, store, auth, userRepository, urlShortenerRepository)

	log.Println(app.Listen(":3000"))

}
