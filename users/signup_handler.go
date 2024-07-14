package users

import (
	"crypto/rand"
	"fmt"
	"net/smtp"

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
			log.Infof("sign-up: %s failed: %v", emailAddress, err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if isExists {
			log.Infof("sign-up: %s failed: email is already registered", emailAddress)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		publicId := uuid.New()
		passwordHash, err := hashPassword(password)
		if err != nil {
			log.Infof("sign-up: %s failed: error on password hashing: %v", emailAddress, err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err = userRepository.StoreUser(publicId, emailAddress, passwordHash); err != nil {
			log.Infof("sign-up: %s failed: %v", emailAddress, err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err = sendVerificationEmail(publicId, emailAddress); err != nil {
			log.Infof("sign-up: %s problem: couldn't send email validation email: %v", emailAddress, err)
			return c.SendStatus(fiber.StatusCreated)
		}

		return c.Redirect("/", fiber.StatusFound)
	}

}

func hashPassword(password string) (*string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashString := string(passwordHash)

	return &hashString, nil
}

func sendVerificationEmail(publicId uuid.UUID, emailAddress string) error {
	host := "localhost:1025"
	from := "noreply@invisibleurl.net"
	to := emailAddress
	subject := "InvisibleURL.Net - Activate your email address"
	body := "Test body " + publicId.String()

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", "", "", "localhost")

	err := smtp.SendMail(host, auth, from, []string{to}, []byte(msg))

	return err
}

func generateActivationTicket(userId int64, userRepository *UserRepository) (*string, error) {

	attempt := 10

	for attempt > 0 {
		token := generateToken()

		if !userRepository.ContainsActivationTicket(userId, token) {
			return token, nil
		}
	}

	return nil, fmt.Errorf("error in generating activation ticket: attempt limit reached")
}

func generateToken() (*string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	token := make([]byte, 50)

	_, err := rand.Read(token)
	if err != nil {
		return nil, fmt.Errorf("error on token generation")
	}

}
