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
		recaptchaSite := environment.RECAPTCHA_SITE
		return c.Render("users/sign-up", fiber.Map{
			"recaptchaSite": recaptchaSite,
		}, "layouts/user")
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
		confirmPassword := c.FormValue("confirmPassword")
		captchaResponse := c.FormValue("g-recaptcha-response")
		remoteIP := c.Context().RemoteIP()

		// Check if this is an HTMX request
		isHtmx := c.Get("HX-Request") == "true"

		// Helper function to return validation errors
		returnValidationError := func(message string) error {
			if isHtmx {
				return c.Status(200).SendString(fmt.Sprintf(`<div id="validation-errors" class="alert alert-danger">%s</div>`, message))
			} else {
				return c.Redirect("/user/sign-up?error="+url.QueryEscape(message), fiber.StatusSeeOther)
			}
		}

		// Check captcha
		if len(captchaResponse) == 0 {
			return returnValidationError("Captcha verification is required")
		}

		// Verify captcha
		if err := verifyCaptcha(captchaResponse); err != nil {
			return returnValidationError("Captcha verification failed: " + err.Error())
		}

		// Validate email
		if err := validateEmail(emailAddress); err != nil {
			return returnValidationError(err.Error())
		}

		// Validate password
		if err := validatePassword(password); err != nil {
			return returnValidationError(err.Error())
		}

		// Check if passwords match
		if err := validateConfirmPassword(password, confirmPassword); err != nil {
			return returnValidationError(err.Error())
		}

		// Check if user already exists
		isExists, err := userRepository.Is_Exists(emailAddress)
		if err != nil {
			log.Infof("sign-up: %s failed: %v", emailAddress, err)
			return returnValidationError("Registration error: please try again later")
		}

		if isExists {
			log.Infof("sign-up: %s failed: email is already registered", emailAddress)
			return returnValidationError("This email address is already registered")
		}

		// Generate public ID and hash password
		publicId := uuid.New()
		passwordHash, err := hashPassword(password)
		if err != nil {
			log.Infof("sign-up: %s failed: error on password hashing: %v", emailAddress, err)
			return returnValidationError("Registration error: could not process password")
		}

		// Store user in database
		var userId int64
		if userId, err = userRepository.StoreUser(publicId, emailAddress, passwordHash); err != nil {
			log.Infof("sign-up: %s failed: %v", emailAddress, err)
			return returnValidationError("Registration error: could not create user")
		}

		// Generate activation ticket
		var activationTicket *string
		if activationTicket, err = generateActivationTicket(userId, userRepository); err != nil {
			log.Infof("sign-up: %s failed: error on creating activation ticket: %v", emailAddress, err)
			return returnValidationError("Registration error: could not create activation ticket")
		}

		// Send verification email
		if err = sendVerificationEmail(emailAddress, *activationTicket); err != nil {
			log.Infof("sign-up: %s problem: couldn't send email validation email: %v", emailAddress, err)
			return returnValidationError("Account created but activation email could not be sent. Please contact support.")
		}

		// Log the registration event
		auditLogService.LogEvent(auditlog.REGISTRATION, userId, remoteIP)

		// Success message and redirect
		if isHtmx {
			successHTML := `<div id="validation-errors" class="alert alert-success">
				Registration successful! Please check your email to activate your account.
				<script>
					setTimeout(function() {
						window.location.href = "/user/sign-in";
					}, 3000);
				</script>
			</div>`
			return c.Status(200).SendString(successHTML)
		}

		return c.Redirect("/user/sign-in", fiber.StatusFound)
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
