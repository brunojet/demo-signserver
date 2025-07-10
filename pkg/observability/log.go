package observability

import (
	"encoding/json"
	"log"
	"time"
)

func logInternal(level, msg string, fields map[string]interface{}) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     level,
		"msg":       msg,
	}
	for k, v := range fields {
		logEntry[k] = v
	}
	jsonLog, _ := json.Marshal(logEntry)
	log.Println(string(jsonLog))
}

// LogInfo registra logs de informação no padrão observabilidade
func LogInfo(msg string, fields map[string]interface{}) {
	logInternal("info", msg, fields)
}

func LogError(msg string, fields map[string]interface{}) {
	logInternal("error", msg, fields)
}
