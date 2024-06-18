package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func GetSignUpHandler() fiber.Handler {

	return func(c *fiber.Ctx) error {
		return c.Render("users/sign-up", fiber.Map{}, "layouts/user")
	}
}

func PostSignUpHHandler(store *session.Store) fiber.Handler {

	return func(c *fiber.Ctx) error {

		emailAddress := c.FormValue("emailAddress")
		password := c.FormValue("password")
		passwordAgain := c.FormValue("passwordAgain")

		log.Info(emailAddress)

		if password != passwordAgain {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		// password validation: https://github.com/go-playground/validator

		// check if email already exists

		// store email

		// send verification email

		return c.SendStatus(fiber.StatusCreated)
	}

}
