package security

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupApp() *fiber.App {

	app := fiber.New()

	RegisterRateLimitingMiddleware(app)

	app.Get("/diag", func(c *fiber.Ctx) error {
		return c.SendString("Ok")
	})

	return app
}

func TestRateLimiting(t *testing.T) {
	app := setupApp()

	for i := 0; i < 100; i++ {
		request, _ := http.NewRequest(http.MethodGet, "/diag", nil)
		resp, err := app.Test(request)

		assert.NoError(t, err)

		if i < 100 {
			assert.Equal(t, http.StatusOK, resp.StatusCode, fmt.Sprintf("Expected status 200 OK. Got %v", resp.StatusCode))
		}
		if i == 100 {
			assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
		}
	}
}
