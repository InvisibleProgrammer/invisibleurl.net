package users

import (
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"invisibleprogrammer.com/invisibleurl/authenticator"
)

func CallbackHandler(store *session.Store, auth *authenticator.Authenticator, repository *UserRepository) fiber.Handler {

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

		// log.Printf("Subject: %v", strings.Split(idToken.Subject, "|")[1])

		userId := extractUserIdFromSubject(idToken.Subject)

		err = repository.StoreUser(userId)
		if err != nil {
			log.Printf("error in storing user: %v", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)
		session.Set("name", profile["name"])
		session.Set("userId", userId)

		if err := session.Save(); err != nil {
			c.SendString(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.Redirect("/", http.StatusTemporaryRedirect)
	}

}

func extractUserIdFromSubject(subject string) string {
	if strings.Contains(subject, "|") {
		return strings.Split(subject, "|")[1]
	} else {
		return subject
	}
}
