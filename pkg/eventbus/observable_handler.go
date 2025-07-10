package eventbus

import (
	"demo-signserver/pkg/observability"
	"time"

	"github.com/google/uuid"
)

type HandlerCtx struct {
	EventType HandlerName
	TraceID   string
	Start     time.Time
}

type ObservableHandler interface {
	HandlerStart(eventType HandlerName) HandlerCtx
	HandlerSuccess(ctx HandlerCtx)
	HandlerError(ctx HandlerCtx, err error)
	HandlerPanic(ctx HandlerCtx, panicVal any)
	HandlerLog(eventCaller, message string, fields map[string]interface{})
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

func (o *DefaultObservableHandler) HandlerStart(eventType HandlerName) HandlerCtx {
	traceID := generateTraceID()
	start := time.Now()
	if o.Obs != nil {
		observability.LogInfo("eventbus.handler.start", map[string]interface{}{
			"eventType": eventType,
			"traceID":   traceID,
		})
		o.Obs.Inc("eventbus.handler.start", map[string]string{"eventType": string(eventType), "traceID": traceID})
	}
	return HandlerCtx{EventType: eventType, TraceID: traceID, Start: start}
}

func (o *DefaultObservableHandler) HandlerSuccess(ctx HandlerCtx) {
	duration := time.Since(ctx.Start)
	if o.Obs != nil {
		observability.LogInfo("eventbus.handler.success", map[string]interface{}{
			"eventType":   ctx.EventType,
			"traceID":     ctx.TraceID,
			"duration_ms": duration.Milliseconds(),
		})
		o.Obs.Inc("eventbus.handler.success", map[string]string{"eventType": string(ctx.EventType), "traceID": ctx.TraceID})
	}
}

func (o *DefaultObservableHandler) HandlerError(ctx HandlerCtx, err error) {
	duration := time.Since(ctx.Start)
	if o.Obs != nil {
		observability.LogInfo("eventbus.handler.error", map[string]interface{}{
			"eventType":   ctx.EventType,
			"traceID":     ctx.TraceID,
			"error":       err.Error(),
			"duration_ms": duration.Milliseconds(),
		})
		o.Obs.Inc("eventbus.handler.error", map[string]string{"eventType": string(ctx.EventType), "traceID": ctx.TraceID})
	}
}

func (o *DefaultObservableHandler) HandlerPanic(ctx HandlerCtx, panicVal any) {
	duration := time.Since(ctx.Start)
	if o.Obs != nil {
		observability.LogInfo("eventbus.handler.panic", map[string]interface{}{
			"eventType":   ctx.EventType,
			"traceID":     ctx.TraceID,
			"panic":       panicVal,
			"duration_ms": duration.Milliseconds(),
		})
		o.Obs.Inc("eventbus.handler.panic", map[string]string{"eventType": string(ctx.EventType), "traceID": ctx.TraceID})
	}
}

func (o *DefaultObservableHandler) HandlerLog(eventCaller, message string, fields map[string]interface{}) {
	if o.Obs != nil {
		fields["eventCaller"] = eventCaller
		fields["traceID"] = o.TraceID
		observability.LogInfo(message, fields)
	}
}

func generateTraceID() string {
	return uuid.NewString()
}
