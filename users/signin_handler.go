package users

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

func GetSignInHandler() fiber.Handler {

	return func(c *fiber.Ctx) error {
		return c.Render("users/sign-in", fiber.Map{}, "layouts/user")
	}
}

func SignInHandler(store *session.Store, userRepository *UserRepository) fiber.Handler {

	return func(c *fiber.Ctx) error {

		emailAddress := c.FormValue("emailAddress")
		password := c.FormValue("password")

		log.Info(emailAddress)

		err := validateEmail(emailAddress)
		if err != nil {
			c.SendString("invalid email")
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err = validatePassword(password)
		if err != nil {
			c.SendString(err.Error())
			c.SendStatus(fiber.StatusBadRequest)
		}

		user, err := userRepository.Get_User_by_Email(emailAddress)
		if err != nil {
			log.Errorf("error on getting user from db: %v", err.Error())
			c.SendString(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		if !user.Activated {
			log.Errorf("cannot log in %s: user is not activated", emailAddress)
			c.SendString(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		if validPassword, err := checkPassword(user.PasswordHash, password); err != nil || !validPassword {
			log.Errorf("password validation failed: %v", err.Error())
			c.SendString(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

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
		session.Set("publicId", user.PublicId)
		session.Set("emailAddress", user.EmailAddress)

		if err := session.Save(); err != nil {
			c.SendString(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.Redirect("/", fiber.StatusFound)
	}

}

func checkPassword(passwordHash, loginPassword string) (bool, error) {
	log.Infof("Password hash: %s, password: %s", passwordHash, loginPassword)

	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(loginPassword))
	if err != nil {
		return false, err
	}

	return true, nil
}
