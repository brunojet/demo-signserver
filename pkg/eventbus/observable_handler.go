package eventbus

import (
	"demo-signserver/pkg/observability"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type HandlerCtx struct {
	EventType HandlerName
	TraceID   string
	Start     time.Time
}

type ObservableHandler interface {
	GenerateTraceID() string
	HandlerStart(eventType HandlerName) HandlerCtx
	HandlerStartWithTraceId(eventType HandlerName, traceID string) HandlerCtx
	HandlerSuccess(ctx HandlerCtx)
	HandlerError(ctx HandlerCtx, err error)
	HandlerPanic(ctx HandlerCtx, panicVal any)
	HandlerLogInfo(caller, message string, fields map[string]interface{})
	HandlerLogError(caller, method string, err error, fields map[string]interface{})
}

type DefaultObservableHandler struct {
	Obs     *observability.MetricsService
	TraceID string
}

func NewDefaultObservableHandler(obs *observability.MetricsService) *DefaultObservableHandler {
	return &DefaultObservableHandler{
		Obs:     obs,
		TraceID: generateTraceID(),
	}
}

func (o *DefaultObservableHandler) HandlerStartWithTraceId(eventType HandlerName, traceID string) HandlerCtx {
	start := time.Now()
	if o.Obs != nil {
		observability.LogInfo("eventbus.handler.start.log", map[string]interface{}{
			"traceID":   traceID,
			"eventType": eventType,
		})
		o.Obs.Inc("eventbus.handler.start.metric", map[string]string{"eventType": string(eventType), "traceID": traceID})
	}
	return HandlerCtx{EventType: eventType, TraceID: traceID, Start: start}
}

func (o *DefaultObservableHandler) HandlerStart(eventType HandlerName) HandlerCtx {
	traceID := o.GenerateTraceID()
	return o.HandlerStartWithTraceId(eventType, traceID)
}

func (o *DefaultObservableHandler) HandlerSuccess(ctx HandlerCtx) {
	if o.Obs != nil {
		duration := time.Since(ctx.Start)
		observability.LogInfo("eventbus.handler.success", map[string]interface{}{
			"eventType":   ctx.EventType,
			"traceID":     ctx.TraceID,
			"duration_ms": duration.Milliseconds(),
		})
		o.Obs.Inc("eventbus.handler.success", map[string]string{"eventType": string(ctx.EventType), "traceID": ctx.TraceID})
	}
}

func (o *DefaultObservableHandler) HandlerError(ctx HandlerCtx, err error) {
	if o.Obs != nil {
		duration := time.Since(ctx.Start)
		observability.LogError("eventbus.handler.error", map[string]interface{}{
			"eventType":   ctx.EventType,
			"traceID":     ctx.TraceID,
			"duration_ms": duration.Milliseconds(),
			"error":       err.Error(),
		})
		o.Obs.Inc("eventbus.handler.error", map[string]string{"eventType": string(ctx.EventType), "traceID": ctx.TraceID})
	}
}

func (o *DefaultObservableHandler) HandlerPanic(ctx HandlerCtx, panicVal any) {
	if o.Obs != nil {
		duration := time.Since(ctx.Start)
		observability.LogError("eventbus.handler.panic", map[string]interface{}{
			"eventType":   ctx.EventType,
			"traceID":     ctx.TraceID,
			"duration_ms": duration.Milliseconds(),
			"panic":       panicVal,
		})
		o.Obs.Inc("eventbus.handler.panic", map[string]string{"eventType": string(ctx.EventType), "traceID": ctx.TraceID})
	}
}

func (o *DefaultObservableHandler) HandlerLogInfo(caller, message string, fields map[string]interface{}) {
	if o.Obs != nil {
		logFields := map[string]interface{}{
			"caller":  caller,
			"traceID": o.TraceID,
		}
		for k, v := range fields {
			logFields[k] = v
		}
		observability.LogInfo(message, fields)
	}
}

func (o *DefaultObservableHandler) HandlerLogError(caller, method string, err error, fields map[string]interface{}) {
	if o.Obs != nil {
		logFields := map[string]interface{}{
			"caller":  caller,
			"traceID": o.TraceID,
			"method":  method,
			"error":   err.Error(),
		}
		for k, v := range fields {
			logFields[k] = v
		}
		observability.LogError(fmt.Sprintf("falha em %s", method), logFields)
	}
}

func generateTraceID() string {
	return uuid.NewString()
}

func (o *DefaultObservableHandler) GenerateTraceID() string {
	return generateTraceID()
}
