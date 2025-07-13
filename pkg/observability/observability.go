package observability

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

type HandlerFunc[T any] func() T

type ObservabilityMiddleware[T any] struct {
	ctx       context.Context
	metrics   *MetricsService
	logger    *log.Logger
	requestID string
}

func NewObservabilityMiddleware[T any](ctx context.Context) *ObservabilityMiddleware[T] {
	sink := SinkFromContext(ctx)
	metrics := NewMetricsService(sink)
	logger := LoggerFromContext(ctx)
	requestID := RequestIDFromContext(ctx)

	if requestID == "" {
		requestID = "unknown-request-id"
	}

	return &ObservabilityMiddleware[T]{
		ctx:       ctx,
		metrics:   metrics,
		logger:    logger,
		requestID: requestID,
	}
}

type logEntry struct {
	Timestamp  string                 `json:"timestamp"`
	Level      string                 `json:"level"`
	RequestID  string                 `json:"requestID"`
	Caller     string                 `json:"caller"`
	Msg        string                 `json:"msg"`
	DurationMs *int64                 `json:"duration_ms,omitempty"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

func (o *ObservabilityMiddleware[T]) log(level, caller, msg string, durationMs *int64, extra map[string]interface{}) {
	if o.logger != nil {
		entry := logEntry{
			Timestamp:  time.Now().Format(time.RFC3339),
			Level:      level,
			RequestID:  o.requestID,
			Caller:     caller,
			Msg:        msg,
			DurationMs: durationMs,
		}
		if len(extra) > 0 {
			entry.Extra = extra
		}
		jsonLog, _ := json.Marshal(entry)
		o.logger.Println(string(jsonLog))
	}
}

func (o *ObservabilityMiddleware[T]) LogStart(level string, caller string, msg string, extra map[string]interface{}) {
	o.log(level, caller, msg, nil, extra)
}

func (o *ObservabilityMiddleware[T]) LogEnd(level string, caller string, msg string, duration time.Duration, extra map[string]interface{}) {
	durationMs := duration.Milliseconds()
	o.log(level, caller, msg, &durationMs, extra)
}

func (o *ObservabilityMiddleware[T]) MetricsInc(name string, tags map[string]string) {
	if o.metrics != nil {
		o.metrics.Inc(name, tags)
	}
}

func (o *ObservabilityMiddleware[T]) Do(caller string, handler HandlerFunc[T]) (result T) {
	start := time.Now()

	o.LogStart("info", caller, "Requisição iniciada", nil)

	o.MetricsInc("http_client_count", map[string]string{"caller": caller, "requestID": o.requestID})

	span, ctxWithSpan := StartSpan(o.ctx, caller, o.requestID)

	defer func() {
		if span != nil {
			span.End()
		}
		duration := time.Since(start)

		if r := recover(); r != nil {
			o.log("error", caller, "Pânico durante execução do handler", nil, map[string]interface{}{
				"panic": r,
			})
		}

		o.LogEnd("info", caller, "Requisição finalizada", duration, nil)
	}()

	result = handlerWithContext(handler, ctxWithSpan)
	return result
}

func handlerWithContext[T any](handler HandlerFunc[T], _ context.Context) T {
	return handler()
}
