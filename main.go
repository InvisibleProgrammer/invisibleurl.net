package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	shortUrl := "blog"
	longUrl := "https://invisibleprogrammer.com"

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("The best URL Shortener ever")
	})

	app.Get("/:shortUrl", func(c *fiber.Ctx) error {
		short := c.Params("shortUrl")

		if short == shortUrl {
			return c.Redirect(longUrl, fiber.StatusFound)
		}

		return c.SendString("waaat")
	})

	app.Listen(":3000")
}
