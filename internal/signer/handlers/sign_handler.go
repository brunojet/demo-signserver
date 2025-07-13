package handlers

import (
	"context"
	"demo-signserver/internal/config"
	"demo-signserver/internal/repository/domain"
	"demo-signserver/internal/repository/repositories"
	"demo-signserver/internal/signer/adapters"
	"demo-signserver/pkg/eventbus"
	"fmt"
)

// Evento: sign_process
func SignProcessHandler(bus *eventbus.EventBus) eventbus.Handler {
	methods := config.GetSignServerMethods()
	return func(ctx context.Context, event any) error {
		request, ok := event.(*domain.SignRequest)

		if !ok {
			return fmt.Errorf("event type mismatch: %v", event)
		}

		var (
			err  error = nil
			step       = "init"
		)

		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic recovered: %v", r)
			}

			if err != nil {
				request.SetSignerStatus(domain.SignerStatusSigningFailed, &domain.SignerError{
					Location: fmt.Sprintf("SignProcessHandler step: %s", step),
					Message:  err.Error(),
				})
			} else {
				request.SetSignerStatus(domain.SignerStatusUploaded, nil)
			}
			repository := repositories.NewRequestRepository()
			repository.UpdateRequest(request.ID, request)
		}()

		request.SetSignerStatus(domain.SignerStatusSigning, nil)

		step = "new_positivo_signer"
		externalSigner, err := adapters.NewPositivoSigner(*request.SignerProfileId)

		if err != nil {
			return err
		}

		storage := methods.NewStorageService()

		srcPath := storage.GetWorkFilePath(request.UnsignedFile)

		step = "start_sign"
		ID, err := externalSigner.StartSign(srcPath)

		if err != nil {
			return err
		}

		dstPath := storage.GetWorkFilePath(request.SignedFile)

		step = "wait_signature"
		err = externalSigner.WaitSignature(ID, dstPath)

		if err != nil {
			return err
		}

		step = "publish_storage_upload"
		err = bus.PublishWithContext(ctx, "storage_upload", request)

		if err != nil {
			return err
		}

		return nil
	}
}
