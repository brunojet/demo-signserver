package handlers

import (
	"context"
	"demo-signserver/pkg/eventbus"
)

type SignProcessEvent struct {
	ID   string // ID do pedido de assinatura
	File string // Caminho do arquivo a ser assinado
}

// Evento: sign_process
func SignProcessHandler(bus *eventbus.EventBus) eventbus.Handler {
	return func(ctx context.Context, event any) {
		// TODO: processar assinatura, publicar evento de upload assinado
		// bus.Publish("signed_upload", ...)
	}
}
