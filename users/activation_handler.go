package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	auditlog "invisibleprogrammer.com/invisibleurl/audit_log"
)

func ActivationHandler(store *session.Store, userRepository *UserRepository, auditLogService *auditlog.AuditLogService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		activationTicket := c.Params("activationTicket")
		remoteIP := c.Context().RemoteIP()

		if activationTicket == "" {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		var userId int64
		var err error
		if userId, err = userRepository.Get_UserId_by_ActivationTicket(activationTicket); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err := userRepository.Activate_User(userId); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		auditLogService.LogEvent(auditlog.EMAIL_ACTIVATION, userId, remoteIP)
		return c.Redirect("/", fiber.StatusFound)
	}
}
