package handlers

import (
	"context"
	"errors"
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
	return func(ctx context.Context, event any) {
		evt, ok := event.(UploadEvent)
		if !ok {
			bus.ObsLogInc("upload.receive.handler", map[string]interface{}{"error": "event type mismatch"}, "upload.receive.error", nil)
			return
		}

		repository := repositories.NewRequestRepository()
		storage := storages.NewStorageService()

		ID := filepath.Base(evt.Key)

		request, err := repository.GetRequestByID(ID)

		if err != nil {
			bus.ObsLogInc("upload.receive.handler", map[string]interface{}{
				"request_id": ID,
				"error":      "Request not found",
			}, "upload.receive.error", nil)
			return
		}

		if (request.UnsignedFile.BucketName != evt.Bucket) ||
			(request.UnsignedFile.ObjectKey != evt.Key) {
			bus.ObsLogInc("upload.receive.handler", map[string]interface{}{
				"request_id": ID,
				"error":      "bucket, key, sha256 or size mismatch",
			}, "upload.receive.error", nil)
			return
		}

		updateRequest := &domain.SignRequest{
			History: request.History,
		}

		defer func() {
			if updateRequest.SignerStatus == nil {
				if err == nil || err.Error() == "" {
					err = errors.New("erro desconhecido")
				}
				updateRequest.SetSignerStatus(domain.SignerStatusSigningFailed, &domain.SignerError{
					Code:    "999",
					Message: err.Error(),
				})
			}
			if err != nil {
				bus.ObsLogInc(
					"upload.receive.handler",
					map[string]interface{}{
						"request_id": ID,
						"status":     updateRequest.SignerStatus,
						"error":      updateRequest.GetLastError(),
					},
					"upload_receive_status",
					nil,
				)
			} else {
				bus.ObsLogInc(
					"upload.receive.handler",
					map[string]interface{}{
						"request_id": ID,
						"status":     updateRequest.SignerStatus,
					},
					"upload_receive_status",
					nil,
				)
			}

			repository.UpdateRequest(ID, updateRequest)
		}()

		err = storage.DownloadFileFromS3(evt.Key)

		if err != nil {
			updateRequest.SetSignerStatus(domain.SignerStatusSigningFailed, &domain.SignerError{
				Code:    "001",
				Message: "Não foi possível baixar o arquivo: " + err.Error(),
			})
			return
		}

		updateRequest.SetUnsignedBucketInfoShaAndSize(evt.ETag, evt.Size)
		updateRequest.SetSignerStatus(domain.SignerStatusUploaded, nil)

		err = bus.Publish("sign_process", SignProcessEvent{
			ID:   ID,
			File: evt.Key,
		})

		if err != nil {
			updateRequest.SetSignerStatus(domain.SignerStatusSigningFailed, &domain.SignerError{
				Code:    "002",
				Message: "Erro ao publicar evento de assinatura: " + err.Error(),
			})
			return
		}
	}
}
