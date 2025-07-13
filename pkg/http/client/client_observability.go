package http_client

import (
	"context"
	"demo-signserver/pkg/observability"
	"encoding/json"
	"log"
	"time"
)

type UploadHandlerFunc func() StatusCode

func logAsync(log *log.Logger, logEntry any) {
	if log != nil {
		jsonLog, _ := json.Marshal(logEntry)
		log.Println(string(jsonLog))
	}
}

func ObservabilityHttpClientMiddleware(ctx context.Context, caller string, handler UploadHandlerFunc) (statusCode StatusCode) {
	start := time.Now()
	statusCode = StatusCode(500)
	sink := observability.SinkFromContext(ctx)
	metrics := observability.NewMetricsService(sink)
	logger := observability.LoggerFromContext(ctx)
	requestID := observability.RequestIDFromContext(ctx)
	if requestID == "" {
		requestID = "unknown-request-id"
	}

	logAsync(logger, map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     "info",
		"requestID": requestID,
		"caller":    caller,
		"msg":       "Requisição iniciada",
	})

	if metrics != nil {
		metrics.Inc("http_client_count", map[string]string{"caller": caller, "requestID": requestID})
	}

	span, ctxWithSpan := observability.StartSpan(ctx, "http_client", requestID)

	defer func() {
		if span != nil {
			span.End()
		}
		duration := time.Since(start)

		logAsync(logger, map[string]interface{}{
			"timestamp":   time.Now().Format(time.RFC3339),
			"level":       "info",
			"requestID":   requestID,
			"caller":      caller,
			"msg":         "Requisição finalizada",
			"statusCode":  statusCode,
			"duration_ms": duration.Milliseconds(),
		})
	}()

	// Handler executado com contexto do span e tratamento de pânico
	defer func() {
		if r := recover(); r != nil {
			logAsync(logger, map[string]interface{}{
				"timestamp": time.Now().Format(time.RFC3339),
				"level":     "error",
				"requestID": requestID,
				"caller":    caller,
				"msg":       "Pânico durante execução do handler",
				"panic":     r,
			})
			statusCode = StatusCode(500)
		}
	}()

	statusCode = handlerWithContext(handler, ctxWithSpan)
	return statusCode
}

func handlerWithContext(handler UploadHandlerFunc, _ context.Context) StatusCode {
	return handler()
}
