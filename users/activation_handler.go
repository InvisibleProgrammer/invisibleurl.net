package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	auditlog "invisibleprogrammer.com/invisibleurl/audit_log"
)

func ActivationHandler(store *session.Store, userRepository *UserRepository, auditLogService *auditlog.AuditLogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		activationTicket := c.Params("activationTicket")
		remoteIP := c.Context().RemoteIP()

		if activationTicket == "" {
			log.Infof("Couldn't activate the user. Activation ticket is empty")
			return c.Render("users/activate-fail", fiber.Map{}, "layouts/user")
		}

		var userId int64
		var err error
		if userId, err = userRepository.Get_UserId_by_ActivationTicket(activationTicket); err != nil || userId == 0 {
			log.Infof("Couldn't activate the user. Didn't find a valid activation ticket %s", activationTicket)
			return c.Render("users/activate-fail", fiber.Map{}, "layouts/user")
		}

		if err := userRepository.Activate_User(userId); err != nil {
			log.Infof("Couldn't activate the user. Activation failed")
			return c.Render("users/activate-fail", fiber.Map{}, "layouts/user")
		}

		auditLogService.LogEvent(auditlog.EMAIL_ACTIVATION, userId, remoteIP)
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
