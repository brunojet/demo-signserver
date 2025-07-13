package handlers

import (
	"context"
	"fmt"

	"demo-signserver/internal/config"
	"demo-signserver/internal/repository/domain"
	"demo-signserver/internal/repository/repositories"
	"demo-signserver/pkg/eventbus"
	"demo-signserver/pkg/storage"
)

type StorageDownloadEvent struct {
	storage.FileInfo
}

func StorageDownloadHandler(bus *eventbus.EventBus) eventbus.Handler {
	methods := config.GetSignServerMethods()
	return func(ctx context.Context, event any) error {
		var (
			err  error = nil
			step       = "init"
		)

		request, ok := event.(*domain.SignRequest)

		if !ok {
			return fmt.Errorf("event type mismatch: %v", event)
		}

		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic recovered: %v", r)
			}

			if err != nil {
				request.SetSignerStatus(domain.SignerStatusSigningFailed, &domain.SignerError{
					Location: fmt.Sprintf("StorageDownloadHandler step: %s", step),
					Message:  err.Error(),
				})
			} else {
				request.SetSignerStatus(domain.SignerStatusUploaded, nil)
			}
			repository := repositories.NewRequestRepository()
			repository.UpdateRequest(request.ID, request)
		}()

		storage := methods.NewStorageService()

		step = "download_file_from_s3"
		err = storage.DownloadFileFromS3(request.UnsignedFile)
		if err != nil {
			return err
		}

		step = "publish_sign_process"
		err = bus.PublishWithContext(ctx, "sign_process", request)

		if err != nil {
			return err
		}

		return nil
	}
}
