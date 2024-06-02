package users

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"invisibleprogrammer.com/invisibleurl/authenticator"
	env "invisibleprogrammer.com/invisibleurl/environment"
)

func LogoutHandler(store *session.Store, auth *authenticator.Authenticator) fiber.Handler {

	return func(c *fiber.Ctx) error {

		logoutUrl, err := url.Parse("https://" + os.Getenv(env.AUTH0_DOMAIN) + "/v2/logout")
		if err != nil {
			c.SendString(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		scheme := "http"
		if c.App().Server().TLSConfig != nil {
			scheme = "https"
		}

		returnTo, err := url.Parse(scheme + "://" + string(c.Request().Host()))
		if err != nil {
			c.SendString(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		parameters := url.Values{}
		parameters.Add("returnTo", returnTo.String())
		parameters.Add("client_id", os.Getenv(env.AUTH0_CLIENT_ID))
		logoutUrl.RawQuery = parameters.Encode()

		session, err := store.Get(c)
		if err != nil {
			log.Fatalf("Couldn't receive sesion: %v", err)
		}

		session.Delete("access_token")
		session.Delete("profile")
		session.Delete("name")
		session.Delete("userId")

		if err := session.Save(); err != nil {
			c.SendString(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.Redirect(logoutUrl.String(), http.StatusTemporaryRedirect)
	}

}
