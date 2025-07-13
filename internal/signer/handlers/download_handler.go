package handlers

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

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

		evt, ok := event.(StorageDownloadEvent)

		if !ok {
			return fmt.Errorf("event type mismatch: %v", event)
		}

		repository := repositories.NewRequestRepository()

		ID := filepath.Base(evt.FilePath)

		request, err := repository.GetRequestByID(ID)

		if err != nil {
			return fmt.Errorf("error getting request by ID %s: %w", ID, err)
		}

		defer func() {
			if r := recover(); r != nil {
				step = "panic_recovery"
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
			repository.UpdateRequest(request.ID, request)
		}()

		storage := methods.NewStorageService()

		err = storage.DownloadFileFromS3(request.UnsignedFile)

		if err != nil {
			step = "download_file_from_s3"
			return err
		}

		log.Printf("[StorageDownloadHandler] File downloaded successfully: %s\n", request.UnsignedFile.FilePath)
		err = bus.PublishWithContext(ctx, "sign_process", request)
		log.Printf("[StorageDownloadHandler] Event published: sign_process for request %s\n", request.ID)

		if err != nil {
			step = "publish_sign_process"
			return err
		}

		return nil
	}
}
