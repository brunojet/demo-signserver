package observability

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ObservableHandlerFunc é uma função de handler que recebe contexto já enriquecido
// com requestID e pode acessar helpers de log/metrics/tracing.
type ObservableHandlerFunc func(c *gin.Context, obs *Observability)

type Observability struct {
	RequestID string
	Logger    *log.Logger // Pode ser trocado por interface customizada
	// Futuro: métricas, tracing, etc
}

// Middleware que injeta observabilidade e executa o handler de negócio
var defaultLogger = log.Default()

func ObservableMiddleware(handler ObservableHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Writer.Header().Set("X-Request-ID", requestID)

		obs := &Observability{
			RequestID: requestID,
			Logger:    defaultLogger,
		}

		start := time.Now()
		logEntry := map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"level":     "info",
			"requestID": requestID,
			"msg":       "Request recebida",
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
		}
		jsonLog, _ := json.Marshal(logEntry)
		obs.Logger.Println(string(jsonLog))

		// --- Métricas: incrementar contador ---
		IncRequestCount(c.FullPath())

		// --- Tracing: iniciar span ---
		span, _ := StartSpan(c.Request.Context(), "http_request", requestID)
		defer span.End()

		handler(c, obs)

		duration := time.Since(start)
		logEntry = map[string]interface{}{
			"timestamp":   time.Now().Format(time.RFC3339),
			"level":       "info",
			"requestID":   requestID,
			"msg":         "Request finalizada",
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status":      c.Writer.Status(),
			"duration_ms": duration.Milliseconds(),
		}
		jsonLog, _ = json.Marshal(logEntry)
		obs.Logger.Println(string(jsonLog))
	}
}

// Exemplo de uso:
// r.POST("/endpoint", observability.ObservableMiddleware(func(c *gin.Context, obs *observability.Observability) {
//     // handler de negócio
// }))
