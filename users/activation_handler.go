package users

import (
	"net/smtp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	auditlog "invisibleprogrammer.com/invisibleurl/audit_log"
	"invisibleprogrammer.com/invisibleurl/environment"
)

func ActivationHandler(store *session.Store, userRepository *UserRepository, auditLogService *auditlog.AuditLogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		activationTicket := c.Params("activationTicket")
		remoteIP := c.Context().RemoteIP()

		if activationTicket == "" {
			log.Infof("Couldn't activate the user. Activation ticket is empty")
			return c.Render("users/activate-fail", fiber.Map{}, "layouts/user")
		}

		var user *User
		var err error
		if user, err = userRepository.Get_User_by_ActivationTicket(activationTicket); err != nil || user == nil {
			log.Infof("Couldn't activate the user. Didn't find a valid activation ticket %s", activationTicket)
			return c.Render("users/activate-fail", fiber.Map{}, "layouts/user")
		}

		if err := userRepository.Activate_User(user.Id); err != nil {
			log.Infof("Couldn't activate the user. Activation failed")
			return c.Render("users/activate-fail", fiber.Map{}, "layouts/user")
		}

		if err = sendSuccessfulActivationEmail(user.EmailAddress); err != nil {
			log.Errorf("Couldn't send successful notification activation for the email %s", user.EmailAddress)
		}

		auditLogService.LogEvent(auditlog.EMAIL_ACTIVATION, user.Id, remoteIP)
		return c.Redirect("/user/activate-successful", fiber.StatusFound)
	}
}

func ActivationSuccessfulHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("users/activate", fiber.Map{}, "layouts/user")
	}
}

func ActivationFailureHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("users/activate-fail", fiber.Map{}, "layouts/user")
	}
}

func sendSuccessfulActivationEmail(emailAddress string) error {
	to := emailAddress
	subject := "InvisibleURL.Net - Successful activation"
	body := "Your email is activated successfully"

	msg := "From: " + environment.EMAIL_FROM + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", environment.EMAIL_FROM, environment.SMTP_PASSWWORD, environment.SMTP_HOST)

	err := smtp.SendMail(environment.SMTP_HOST+":"+environment.SMTP_PORT, auth, environment.EMAIL_FROM, []string{to}, []byte(msg))

	return err
}
