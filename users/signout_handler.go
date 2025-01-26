package users

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	auditlog "invisibleprogrammer.com/invisibleurl/audit_log"
)

func SignOutHandler(store *session.Store, userRepository *UserRepository, auditLogService *auditlog.AuditLogService) fiber.Handler {

	return func(c *fiber.Ctx) error {

		remoteIP := c.Context().RemoteIP()

		session, err := store.Get(c)
		if err != nil {
			log.Fatalf("Couldn't receive sesion: %v", err)
		}

		publicId := session.Get("publicId")
		user, err := userRepository.Get_UserId_by_PublicId(publicId.(string))
		if err != nil {
			errorMessage := fmt.Sprintf("Cannot get user by public id: %s", err)
			log.Print(errorMessage)

			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": errorMessage,
			})
		}

		session.Delete("state")
		session.Delete("userId")
		session.Delete("emailAddress")
		session.Delete("publicId")

		if err := session.Save(); err != nil {
			c.SendString(err.Error())
			return c.SendStatus(http.StatusInternalServerError)
		}

		auditLogService.LogEvent(auditlog.LOGOUT, user.Id, remoteIP)
		return c.Redirect("/", fiber.StatusFound)
	}

}
