package handlers

import (
	"context"
	"fmt"

	"demo-signserver/internal/config"
	"demo-signserver/internal/repository/domain"
	"demo-signserver/internal/repository/repositories"
	"demo-signserver/pkg/eventbus"
)

func StorageUploadHandler(bus *eventbus.EventBus) eventbus.Handler {
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
					Location: fmt.Sprintf("StorageUploadHandler step: %s", step),
					Message:  err.Error(),
				})
			} else {
				request.SetSignerStatus(domain.SignerStatusSignedAvailable, nil)
			}
			repository := repositories.NewRequestRepository()

			repository.UpdateRequest(request.ID, request)
		}()

		step = "upload_to_s3"
		storage := methods.NewStorageService(ctx)
		err = storage.UploadToS3(request.SignedFile)

		if err != nil {
			return err
		}

		return nil
	}
}
