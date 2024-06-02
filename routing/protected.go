package routing

import "github.com/gofiber/fiber/v2"

/*
For diagnostic/test only. It checks if protection is working for pages
that should be visible only for logged-in users.
*/
func ProtectedHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("protected", fiber.Map{
			"Title": "Hello, World. You shouldn't see that unless you logged in!",
		}, "layouts/main")
	}
}
