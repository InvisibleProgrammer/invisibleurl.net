package users

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"invisibleprogrammer.com/invisibleurl/environment"
)

func GetSignInHandler() fiber.Handler {

	return func(c *fiber.Ctx) error {
		recaptchaSite := environment.RECAPTCHA_SITE
		needCaptcha := c.QueryBool("needCaptcha", false)

		return c.Render("users/sign-in", fiber.Map{
			"needCaptcha":   needCaptcha,
			"recaptchaSite": recaptchaSite,
		}, "layouts/user")
	}
}

func SignInHandler(store *session.Store, userRepository *UserRepository) fiber.Handler {

	return func(c *fiber.Ctx) error {

		emailAddress := c.FormValue("emailAddress")
		password := c.FormValue("password")
		captchaResponse := c.FormValue("g-recaptcha-response")
		haveCaptchaResponse := len(captchaResponse) > 0

		if err := verifyCaptcha(captchaResponse); err != nil {
			c.SendString(err.Error())
			c.SendStatus(fiber.StatusBadRequest)
		}

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
			return c.SendStatus(http.StatusInternalServerError)
		}

		if validPassword, err := checkPassword(user.PasswordHash, password); err != nil || !validPassword {
			log.Errorf("password validation failed: %v", err.Error())
			c.SendString(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		remoteIP := c.Context().RemoteIP()
		if remoteIP == nil {
			message := "IP address check failed: remote IP is not provided"
			log.Error(message)
			return c.SendStatus(http.StatusInternalServerError)
		}

		knownIP, err := userRepository.Is_Known_IP(user.Id, remoteIP)
		if err != nil {
			log.Errorf("IP address check failed: %v", err.Error())
			c.SendString(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		if !knownIP && !haveCaptchaResponse {
			return c.Redirect("/user/sign-in?needCaptcha=true")
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

		if err := userRepository.StoreNewIP(user.Id, remoteIP); err != nil {
			log.Warnf("Couldn't store new IP location for user: %d, remoteIP: %s", user.Id, remoteIP.String())
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
