package users

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/smtp"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	auditlog "invisibleprogrammer.com/invisibleurl/audit_log"
	"invisibleprogrammer.com/invisibleurl/environment"
)

func GetSignUpHandler() fiber.Handler {

	return func(c *fiber.Ctx) error {

		return c.Render("users/sign-up", fiber.Map{}, "layouts/user")
	}
}

type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes,omitempty"`
}

func PostSignUpHandler(store *session.Store, userRepository *UserRepository, auditLogService *auditlog.AuditLogService) fiber.Handler {

	return func(c *fiber.Ctx) error {

		emailAddress := c.FormValue("emailAddress")
		password := c.FormValue("password")
		passwordAgain := c.FormValue("passwordAgain")
		captchaResponse := c.FormValue("g-recaptcha-response")
		haveCaptchaResponse := len(captchaResponse) > 0
		remoteIP := c.Context().RemoteIP()

		if !haveCaptchaResponse {
			return c.Render("user/sign-up?needCaptcha=true", fiber.Map{}, "layouts/user")
		}

		if err := verifyCaptcha(captchaResponse); err != nil {
			c.SendString(err.Error())
			c.SendStatus(fiber.StatusBadRequest)
		}

		log.Info(emailAddress)

		if err := validateEmail(emailAddress); err != nil {
			c.SendString("invalid email")
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err := validateConfirmPassword(password, passwordAgain); err != nil {
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

		var userId int64
		if userId, err = userRepository.StoreUser(publicId, emailAddress, passwordHash); err != nil {
			log.Infof("sign-up: %s failed: %v", emailAddress, err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		var activationTicket *string
		if activationTicket, err = generateActivationTicket(userId, userRepository); err != nil {
			log.Infof("sign-up: %s failed: error on creating activation ticket: %v", emailAddress, err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err = sendVerificationEmail(emailAddress, *activationTicket); err != nil {
			log.Infof("sign-up: %s problem: couldn't send email validation email: %v", emailAddress, err)
			return c.SendStatus(fiber.StatusCreated)
		}

		auditLogService.LogEvent(auditlog.REGISTRATION, userId, remoteIP)
		return c.Redirect("/", fiber.StatusFound)
	}

}

func verifyCaptcha(captchaResponse string) error {

	resp, err := http.PostForm(environment.RECAPTCHA_VERIFY_URL, url.Values{
		"secret":   {environment.RECAPTCHA_SECRET},
		"response": {captchaResponse},
	})

	if err != nil {
		return fmt.Errorf("failed to verify reCAPTCHA")
	}
	defer resp.Body.Close()

	var recaptchaResponse RecaptchaResponse
	if err := json.NewDecoder(resp.Body).Decode(&recaptchaResponse); err != nil {
		return fmt.Errorf("failed to parse reCAPTCHA response")
	}

	log.Infof("Recaptcha success: %t", recaptchaResponse.Success)

	if !recaptchaResponse.Success {
		return fmt.Errorf("reCAPTCHA verification failed")
	}
	return nil
}

func hashPassword(password string) (*string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashString := string(passwordHash)

	return &hashString, nil
}

func sendVerificationEmail(emailAddress string, activationTicket string) error {
	host := "localhost:1025"
	from := "noreply@invisibleurl.net"
	to := emailAddress
	subject := "InvisibleURL.Net - Activate your email address"
	body := fmt.Sprintf("Please activate: https://localhost:3000/user/activate/%s", activationTicket)

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", "", "", "localhost")

	err := smtp.SendMail(host, auth, from, []string{to}, []byte(msg))

	return err
}

func generateActivationTicket(userId int64, userRepository *UserRepository) (*string, error) {

	for attempt := 10; attempt > 0; attempt-- {
		token, err := generateToken()

		if err != nil {
			return nil, err
		}

		if err = userRepository.StoreActivationTicket(userId, token); err != nil {
			log.Infof("couldn't generate activation ticket for user: %d, attempt: %d", userId, 10-attempt+1)
			continue
		}

		return token, nil
	}

	return nil, fmt.Errorf("error in generating activation ticket: attempt limit reached")
}

func generateToken() (*string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 50)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return nil, err
		}
		result[i] = charset[num.Int64()]
	}

	stringToken := string(result)
	return &stringToken, nil
}
