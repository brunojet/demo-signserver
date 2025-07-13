package http_client

import (
	"context"
	"demo-signserver/pkg/observability"
	"encoding/json"
	"log"
	"time"
)

type HandlerFunc func() StatusCode

type ObservabilityHttpClientMiddleware struct {
	ctx       context.Context
	metrics   *observability.MetricsService
	logger    *log.Logger
	requestID string
}

func NewObservabilityHttpClientMiddleware(ctx context.Context) *ObservabilityHttpClientMiddleware {
	sink := observability.SinkFromContext(ctx)
	metrics := observability.NewMetricsService(sink)
	logger := observability.LoggerFromContext(ctx)
	requestID := observability.RequestIDFromContext(ctx)

	if requestID == "" {
		requestID = "unknown-request-id"
	}

	return &ObservabilityHttpClientMiddleware{
		ctx:       ctx,
		metrics:   metrics,
		logger:    logger,
		requestID: requestID,
	}
}

func (o *ObservabilityHttpClientMiddleware) log(logEntry any) {
	if o.logger != nil {
		jsonLog, _ := json.Marshal(logEntry)
		o.logger.Println(string(jsonLog))
	}
}

func (o *ObservabilityHttpClientMiddleware) do(caller string, handler HandlerFunc) StatusCode {
	start := time.Now()
	statusCode := StatusCode(500)

	o.log(map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     "info",
		"requestID": o.requestID,
		"caller":    caller,
		"msg":       "Requisição iniciada",
	})

	if o.metrics != nil {
		o.metrics.Inc("http_client_count", map[string]string{"caller": caller, "requestID": o.requestID})
	}

	span, ctxWithSpan := observability.StartSpan(o.ctx, "http_client", o.requestID)

	defer func() {
		if span != nil {
			span.End()
		}
		duration := time.Since(start)

		if r := recover(); r != nil {
			o.log(map[string]interface{}{
				"timestamp": time.Now().Format(time.RFC3339),
				"level":     "error",
				"requestID": o.requestID,
				"caller":    caller,
				"msg":       "Pânico durante execução do handler",
				"panic":     r,
			})
			statusCode = StatusCode(500)
		}

		o.log(map[string]interface{}{
			"timestamp":   time.Now().Format(time.RFC3339),
			"level":       "info",
			"requestID":   o.requestID,
			"caller":      caller,
			"msg":         "Requisição finalizada",
			"statusCode":  statusCode,
			"duration_ms": duration.Milliseconds(),
		})
	}()

	statusCode = handlerWithContext(handler, ctxWithSpan)
	return statusCode
}

func handlerWithContext(handler HandlerFunc, _ context.Context) StatusCode {
	return handler()
}
