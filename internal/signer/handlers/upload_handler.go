package handlers

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"demo-signserver/internal/repository/domain"
	"demo-signserver/internal/repository/repositories"
	storages "demo-signserver/internal/storage"
	"demo-signserver/pkg/eventbus"
)

type UploadEvent struct {
	Bucket string
	Key    string
	ETag   string
	Size   int64
}

func UploadReceivedHandler(bus *eventbus.EventBus) eventbus.Handler {
	return func(ctx context.Context, event any) error {
		evt, ok := event.(UploadEvent)

		if !ok {
			return fmt.Errorf("event type mismatch: %v", event)
		}

		repository := repositories.NewRequestRepository()
		storage := storages.NewStorageService()

		ID := filepath.Base(evt.Key)

		request, err := repository.GetRequestByID(ID)

		if err != nil {
			return fmt.Errorf("error getting request by ID %s: %w", ID, err)
		}

		if (request.UnsignedFile.Bucket != evt.Bucket) ||
			(request.UnsignedFile.Key != evt.Key) {
			return fmt.Errorf("bucket or key mismatch for request %s", ID)
		}

		updateRequest := &domain.SignRequest{
			History: request.History,
		}

		defer func() {
			if updateRequest.SignerStatus == nil {
				log.Fatalf("signer status is nil for request %s", ID)
			}
			repository.UpdateRequest(ID, updateRequest)
		}()

		err = storage.DownloadFileFromS3(evt.Key)

		if err != nil {
			updateRequest.SetSignerStatus(domain.SignerStatusSigningFailed, &domain.SignerError{
				Code:    "001",
				Message: "Não foi possível baixar o arquivo: " + err.Error(),
			})
			return err
		}

		updateRequest.SetUnsignedBucketInfo(evt.Bucket, evt.Key, evt.ETag, evt.Size)
		updateRequest.SetSignerStatus(domain.SignerStatusUploaded, nil)

		err = bus.Publish("sign_process", ID)

		if err != nil {
			updateRequest.SetSignerStatus(domain.SignerStatusSigningFailed, &domain.SignerError{
				Code:    "002",
				Message: "Erro ao publicar evento de assinatura: " + err.Error(),
			})
			return err
		}

		return nil
	}
}
