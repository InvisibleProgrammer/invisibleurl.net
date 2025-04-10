package users

import (
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
	auditlog "invisibleprogrammer.com/invisibleurl/audit_log"
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

func SignInHandler(store *session.Store, userRepository *UserRepository, auditlogService *auditlog.AuditLogService) fiber.Handler {

	return func(c *fiber.Ctx) error {
		emailAddress := c.FormValue("emailAddress")
		password := c.FormValue("password")
		captchaResponse := c.FormValue("g-recaptcha-response")
		haveCaptchaResponse := len(captchaResponse) > 0
		remoteIP := c.Context().RemoteIP()

		// Check if this is an HTMX request
		isHtmx := c.Get("HX-Request") == "true"

		// Helper function to return validation errors
		returnValidationError := func(message string) error {
			if isHtmx {
				return c.Status(200).SendString(fmt.Sprintf(`<div id="validation-errors" class="alert alert-danger">%s</div>`, message))
			} else {
				return c.Redirect("/user/sign-in?error="+url.QueryEscape(message), fiber.StatusSeeOther)
			}
		}

		// Verify captcha
		if haveCaptchaResponse {
			if err := verifyCaptcha(captchaResponse); err != nil {
				return returnValidationError("Captcha verification failed: " + err.Error())
			}
		}

		// Validate email
		if err := validateEmail(emailAddress); err != nil {
			return returnValidationError(err.Error())
		}

		// Validate password
		if err := validatePassword(password); err != nil {
			return returnValidationError(err.Error())
		}

		// Fetch user by email
		user, err := userRepository.Get_User_by_Email(emailAddress)
		if err != nil {
			log.Errorf("error on getting user from db: %v", err.Error())
			return returnValidationError("Invalid email or password")
		}

		// Check if user is activated
		if !user.Activated {
			log.Errorf("cannot log in %s: user is not activated", emailAddress)
			return returnValidationError("Your account is not activated. Please check your email for activation instructions.")
		}

		// Validate password hash
		if validPassword, err := checkPassword(user.PasswordHash, password); err != nil || !validPassword {
			log.Errorf("password validation failed: %v", err)
			return returnValidationError("Invalid email or password")
		}

		// Check IP address
		if remoteIP == nil {
			message := "IP address check failed: remote IP is not provided"
			log.Error(message)
			return returnValidationError("Authentication error: IP address check failed")
		}

		// Check if IP is known
		knownIP, err := userRepository.Is_Known_IP(user.Id, remoteIP)
		if err != nil {
			log.Errorf("IP address check failed: %v", err.Error())
			return returnValidationError("Authentication error: IP check failed")
		}

		// Require captcha for unknown IP
		if !knownIP && !haveCaptchaResponse {
			if isHtmx {
				c.Set("HX-Redirect", "/user/sign-in?needCaptcha=true")
				return c.Status(200).SendString("")
			}
			return c.Redirect("/user/sign-in?needCaptcha=true")
		}

		// Generate session state
		state, err := generateRandomState()
		if err != nil {
			return returnValidationError("Authentication error: could not create session")
		}

		// Get session
		session, err := store.Get(c)
		if err != nil {
			log.Errorf("Couldn't receive session: %v", err)
			return returnValidationError("Authentication error: could not create session")
		}

		// Set session data
		session.Set("state", state)
		session.Set("publicId", user.PublicId)
		session.Set("emailAddress", user.EmailAddress)

		// Save session
		if err := session.Save(); err != nil {
			return returnValidationError("Authentication error: could not save session")
		}

		// Store new IP
		if err := userRepository.StoreNewIP(user.Id, remoteIP); err != nil {
			log.Warnf("Couldn't store new IP location for user: %d, remoteIP: %s", user.Id, remoteIP.String())
		}

		// Log audit event
		auditlogService.LogEvent(auditlog.LOGIN, user.Id, remoteIP)

		// Success - redirect
		if isHtmx {
			c.Set("HX-Redirect", "/")
			return c.Status(200).SendString("")
		}
		return c.Redirect("/", fiber.StatusFound)
	}
}

func checkPassword(passwordHash, loginPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(loginPassword))
	if err != nil {
		return false, err
	}

	return true, nil
}
