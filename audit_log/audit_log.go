package auditlog

import (
	"net"
)

type AuditLogService struct {
	db AuditLogRepository
}

type AuditLog interface {
	Log(action int)
	LogWithDetail(action int, detail string)
}

type Action int

const (
	REGISTRATION            = 1
	EMAIL_ACTIVATION        = 2
	LOGIN                   = 3
	LOGOUT                  = 4
	CREATE_SHORT_URL        = 10
	CREATE_CUSTOM_SHORT_URL = 11
	DELETE_SHORTENED_URL    = 12
)

var actionNames = map[Action]string{
	REGISTRATION:            "REGISTRATION",
	EMAIL_ACTIVATION:        "EMAIL_ACTIVATION",
	LOGIN:                   "LOGIN",
	LOGOUT:                  "LOGOUT",
	CREATE_SHORT_URL:        "CREATE_SHORT_URL",
	CREATE_CUSTOM_SHORT_URL: "CREATE_CUSTOM_SHORT_URL",
	DELETE_SHORTENED_URL:    "DELETE_SHORTENED_URL",
}

func NewAuditLogService(repository *AuditLogRepository) *AuditLogService {
	return &AuditLogService{
		db: *repository,
	}
}

func (auditLogService *AuditLogService) LogEvent(event Action, userId int64, remoteIP net.IP) {
	auditLogService.db.LogEvent(int(event), userId, remoteIP)
}
