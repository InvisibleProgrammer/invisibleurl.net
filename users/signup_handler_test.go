package users

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRegisterEndpoint(t *testing.T) {
	app := fiber.New()

	// Define the /diag route
	app.Post("/user/register", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Create a request to pass to our handler
	req := httptest.NewRequest("POST", "/user/register", nil)

	// Perform the request
	resp, err := app.Test(req, -1) // -1 disables the request timeout
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// Check if the status code is 200 OK
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// // Check the response body
	// body := make([]byte, resp.ContentLength)
	// resp.Body.Read(body)
	// assert.Equal(t, "OK", string(body))
}
