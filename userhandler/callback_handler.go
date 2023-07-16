package userhandler

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"invisibleprogrammer.com/invisibleurl/authenticator"
)

func CallbackHandler(store *session.Store, auth *authenticator.Authenticator) fiber.Handler {

	return func(c *fiber.Ctx) error {

		session, err := store.Get(c)
		if err != nil {
			log.Fatalf("Couldn't receive sesion: %v", err)
		}

		if c.Query("state") != session.Get("state") {
			c.SendString(("Invalid state parameter."))
			return c.SendStatus(http.StatusBadRequest)
		}

		token, err := auth.Exchange(c.Context(), c.Query("code"))
		if err != nil {
			c.SendString("Failed to exchange an authorization code for a token.")
			return c.SendStatus(http.StatusUnauthorized)
		}

		idToken, err := auth.VerifyIDToken(c.Context(), token)
		if err != nil {
			c.SendString("Failed to verify ID Token.")
			return c.SendStatus(http.StatusInternalServerError)
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			c.SendString(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)
		session.Set("name", profile["name"])

		if err := session.Save(); err != nil {
			c.SendString(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.Redirect("/", http.StatusTemporaryRedirect)
	}

}
