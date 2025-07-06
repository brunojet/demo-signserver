package signer

import (
	"context"
	"demo-signserver/pkg/eventbus"
)

// Evento: signed_upload
func SignedUploadHandler(bus *eventbus.EventBus) eventbus.Handler {
	return func(ctx context.Context, event any) {
		// TODO: fazer upload do arquivo assinado para o S3/storage
	}
}
