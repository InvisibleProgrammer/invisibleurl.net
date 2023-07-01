package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	shortUrl := "blog"
	longUrl := "https://invisibleprogrammer.com"

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.Render("protected", fiber.Map{
			"Title": "Hello, World. You shouldn't see that unless you logged in!",
		}, "layouts/main")
	})

	app.Get("/:shortUrl", func(c *fiber.Ctx) error {
		short := c.Params("shortUrl")

		if short == shortUrl {
			return c.Redirect(longUrl, fiber.StatusFound)
		}

		return c.SendString("waaat")
	})

	log.Fatal(app.Listen(":3000"))
}
