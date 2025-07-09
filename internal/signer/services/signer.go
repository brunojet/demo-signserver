package services

import (
	"demo-signserver/internal/config"
	"demo-signserver/internal/signer/handlers"
	"demo-signserver/pkg/eventbus"
	message_adapters "demo-signserver/pkg/message/adapters"

	"github.com/aws/aws-lambda-go/events"
)

type SignerService struct {
	MessageQueueAdapter message_adapters.MessageQueueAdapterInterface
	EventBus            *eventbus.EventBus
}

func NewSignerService() *SignerService {
	cfg := config.GetSignServerConfig()
	methods := config.GetSignServerMethods()
	message_queue := methods.MessageQueueAdapter(cfg.StorageBucketName, "unsigned")
	return &SignerService{
		EventBus:            eventbus.NewEventBusWithSink(nil),
		MessageQueueAdapter: message_queue,
	}
}

func (s *SignerService) Start() {
	s.EventBus.Register("upload_received", handlers.UploadReceivedHandler(s.EventBus), 2, 10)
	s.EventBus.Register("sign_process", handlers.SignProcessHandler(s.EventBus), 10, 20)
	s.MessageQueueAdapter.Start(func(event any) {
		s3evt, ok := event.(events.S3Event)
		if ok {
			for _, record := range s3evt.Records {
				s.EventBus.Publish("upload_received", handlers.UploadEvent{
					Bucket: record.S3.Bucket.Name,
					Key:    record.S3.Object.Key,
					ETag:   record.S3.Object.ETag,
					Size:   record.S3.Object.Size,
				})
			}
			return
		}
	})

}
