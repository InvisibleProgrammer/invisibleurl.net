package auditlog

import (
	"net"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"invisibleprogrammer.com/invisibleurl/db"
)

type AuditLogRepository struct {
	db *db.Repository
}

func NewAuditLogRepository(repository *db.Repository) *AuditLogRepository {
	return &AuditLogRepository{
		db: repository,
	}
}

func (repository *AuditLogRepository) LogEvent(eventId int, userId int64, remoteIP net.IP) {
	insertStmnt := `insert into audit_log(user_id, IP_Address, event_id, recorded_at) values (:userId, :ipAddress, :eventId, :recordedAt)`

	parameters := map[string]interface{}{
		"userId":     userId,
		"ipAddress":  remoteIP.String(),
		"eventId":    eventId,
		"recordedAt": time.Now(),
	}

	result, err := repository.db.Db.NamedExec(insertStmnt, parameters)
	if err != nil {
		log.Errorf("error on storing audit log entry: Event: %d - User: %d - Error: %v", eventId, userId, err)
		return
	}

	if affectedRows, err := result.RowsAffected(); err != nil || affectedRows == 0 {
		log.Errorf("error on storing audit log entry: Event: %d - User: %d - Error: %v", eventId, userId, err)
		return
	}
}
