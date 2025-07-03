package observability

import (
	"encoding/json"
	"log"
	"time"
)

// Registra uma operação sensível para auditoria
func AuditLog(requestID, user, action string, payload interface{}) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     "audit",
		"requestID": requestID,
		"user":      user,
		"action":    action,
		"payload":   payload,
	}
	jsonLog, _ := json.Marshal(logEntry)
	log.Println(string(jsonLog))
}
