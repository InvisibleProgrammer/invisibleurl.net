package users

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"invisibleprogrammer.com/invisibleurl/authenticator"
)

func LoginHandler(store *session.Store, auth *authenticator.Authenticator) fiber.Handler {

	return func(c *fiber.Ctx) error {

		state, err := generateRandomState()
		if err != nil {
			c.SendString(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		session, err := store.Get(c)
		if err != nil {
			log.Fatalf("Couldn't receive sesion: %v", err)
		}

		session.Set("state", state)

		if err := session.Save(); err != nil {
			c.SendString(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.Redirect(auth.AuthCodeURL(state), http.StatusTemporaryRedirect)
	}

}

func generateRandomState() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
