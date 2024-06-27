package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetSignUpHandler() fiber.Handler {

	return func(c *fiber.Ctx) error {
		return c.Render("users/sign-up", fiber.Map{}, "layouts/user")
	}
}

func PostSignUpHHandler(store *session.Store, userRepository *UserRepository) fiber.Handler {

	return func(c *fiber.Ctx) error {

		emailAddress := c.FormValue("emailAddress")
		password := c.FormValue("password")
		passwordAgain := c.FormValue("passwordAgain")

		log.Info(emailAddress)

		err := validateEmail(emailAddress)
		if err != nil {
			c.SendString("invalid email")
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err = validatePassword(password, passwordAgain)
		if err != nil {
			c.SendString(err.Error())
			c.SendStatus(fiber.StatusBadRequest)
		}

		if password != passwordAgain {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		isExists, err := userRepository.Is_Exists(emailAddress)
		if err != nil {
			log.Info("sign-up: %s failed: %v", emailAddress, err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if isExists {
			log.Info("sign-up: %s failed: email is already registered", emailAddress)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		publicId := uuid.New()
		passwordHash, err := hashPassword(password)
		if err != nil {
			log.Info("sign-up: %s failed: error on password hashing: %v", emailAddress, err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err = userRepository.StoreUser(publicId, emailAddress, &passwordHash); err != nil {
			log.Info("sign-up: %s failed: %v", emailAddress, err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err = sendVerificationEmail(publicId, emailAddress); err != nil {
			log.Info("sign-up: %s problem: couldn't send email validation email: %v", emailAddress, err)
			return c.SendStatus(fiber.StatusCreated)
		}

		return c.SendStatus(fiber.StatusCreated)
	}

}

func hashPassword(password string) (*string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return passwordHash, nil
}
