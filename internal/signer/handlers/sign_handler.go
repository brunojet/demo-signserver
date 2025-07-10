package handlers

import (
	"context"
	"demo-signserver/pkg/eventbus"
)

type SignProcessEvent string

// Evento: sign_process
func SignProcessHandler(bus *eventbus.EventBus) eventbus.Handler {
	return func(ctx context.Context, event any) error {
		// TODO: processar assinatura, publicar evento de upload assinado
		// bus.Publish("signed_upload", ...)
		return nil
	}
}
