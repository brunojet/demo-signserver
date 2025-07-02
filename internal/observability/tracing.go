package observability

import (
	"context"
	"log"
	"time"
)

type Span struct {
	Name      string
	StartTime time.Time
	RequestID string
}

// Inicia um span de tracing (dummy, só loga)
func StartSpan(ctx context.Context, name string, requestID string) (*Span, context.Context) {
	span := &Span{
		Name:      name,
		StartTime: time.Now(),
		RequestID: requestID,
	}
	log.Printf(`{"timestamp":"%s","level":"info","requestID":"%s","span":"%s","msg":"Span iniciado"}`,
		span.StartTime.Format(time.RFC3339), requestID, name)
	return span, ctx
}

// Finaliza um span de tracing (dummy, só loga)
func (s *Span) End() {
	endTime := time.Now()
	duration := endTime.Sub(s.StartTime)
	log.Printf(`{"timestamp":"%s","level":"info","requestID":"%s","span":"%s","msg":"Span finalizado","duration_ms":%d}`,
		endTime.Format(time.RFC3339), s.RequestID, s.Name, duration.Milliseconds())
}
