package storages

import (
	"demo-signserver/internal/config"
	"demo-signserver/pkg/storage/adapters"
	"log"
	"time"
)

type StorageService struct {
	storage adapters.StorageServiceInterface
}

func NewStorageService() adapters.StorageServiceInterface {
	cfg := config.GetSignServerConfig()
	methods := config.GetSignServerMethods()
	log.Printf("Initializing StorageService with bucket name: %s", cfg.StorageBucketName)
	return &StorageService{
		storage: methods.NewStorageService(cfg.StorageBucketName),
	}
}

func (s *StorageService) GetBucketName() string {
	return s.storage.GetBucketName()
}

func (s *StorageService) GeneratePresignedURL(httpMethod adapters.HttpMethod, key string, expires time.Duration) (string, error) {
	return s.storage.GeneratePresignedURL(httpMethod, key, expires)
}

func (s *StorageService) DownloadFile(key, dest string) error {
	return s.storage.DownloadFile(key, dest)
}

func (s *StorageService) UploadFile(key, src string) error {
	return s.storage.UploadFile(key, src)
}
