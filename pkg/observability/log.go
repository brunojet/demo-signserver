package observability

import (
	"encoding/json"
	"log"
	"time"
)

// LogInfo registra logs de informação no padrão observabilidade
func LogInfo(msg string, fields map[string]interface{}) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     "info",
		"msg":       msg,
	}
	for k, v := range fields {
		logEntry[k] = v
	}
	jsonLog, _ := json.Marshal(logEntry)
	log.Println(string(jsonLog))
}
