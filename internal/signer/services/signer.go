package services

import (
	"demo-signserver/internal/config"
	"demo-signserver/internal/signer/handlers"
	"demo-signserver/pkg/eventbus"
	message_adapters "demo-signserver/pkg/message/adapters"
	"demo-signserver/pkg/storage"

	"github.com/aws/aws-lambda-go/events"
)

type SignerService struct {
	Queue    message_adapters.MessageQueueAdapterInterface
	EventBus *eventbus.EventBus
}

func NewSignerService() *SignerService {
	cfg := config.GetSignServerConfig()
	methods := config.GetSignServerMethods()
	return &SignerService{
		Queue:    methods.NewMessageQueueAdapter(cfg.StorageBucketName, "unsigned"),
		EventBus: methods.NewEventBus(),
	}
}

func (s *SignerService) Start() {
	s.EventBus.Register("storage_download", handlers.StorageDownloadHandler(s.EventBus), 2, 10)
	s.EventBus.Register("sign_process", handlers.SignProcessHandler(s.EventBus), 10, 20)
	s.EventBus.Register("storage_upload", handlers.StorageUploadHandler(s.EventBus), 2, 10)
	s.Queue.Start(func(event any) {
		s3evt, ok := event.(events.S3Event)
		if ok {
			for _, record := range s3evt.Records {
				event := handlers.StorageDownloadEvent{}
				event.StoragePath = record.S3.Bucket.Name
				event.FilePath = record.S3.Object.Key
				event.Hash = record.S3.Object.ETag
				event.Size = record.S3.Object.Size
				event.State = storage.FileStatePending
				s.EventBus.Publish("storage_download", event)
			}
			return
		}
	})
}
