package observability

import (
	"context"
	"log"
)

// traceIDKeyType é um tipo privado para evitar colisão de chave no contexto
// Exportamos apenas as funções utilitárias para manipular o traceID

type traceIDKeyType struct{}

var traceIDKey = traceIDKeyType{}

// ContextWithTraceID retorna um novo contexto com o traceID definido
func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// RequestIDFromContext recupera o traceID do contexto, se existir
func RequestIDFromContext(ctx context.Context) string {
	if v := ctx.Value(traceIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

type loggerKeyType struct{}

var loggerKey = loggerKeyType{}

// ContextWithLogger retorna um novo contexto com o logger definido
func ContextWithLogger(ctx context.Context, logger *log.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// LoggerFromContext recupera o logger do contexto, se existir
func LoggerFromContext(ctx context.Context) *log.Logger {
	if v := ctx.Value(loggerKey); v != nil {
		if l, ok := v.(*log.Logger); ok {
			return l
		}
	}
	return nil
}
