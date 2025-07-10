package storages

import (
	"demo-signserver/internal/config"
	"demo-signserver/pkg/storage/adapters"
	"io"
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

// DownloadFileFromS3 implements adapters.StorageServiceInterface.
func (s *StorageService) DownloadFileFromS3(key string) error {
	return s.storage.DownloadFileFromS3(key)
}

// UploadToS3 implements adapters.StorageServiceInterface.
func (s *StorageService) UploadToS3(key string) error {
	return s.storage.UploadToS3(key)
}

// OpenWorkFile implements adapters.StorageServiceInterface.
func (s *StorageService) OpenWorkFile(key string) (io.ReadWriteCloser, error) {
	return s.storage.OpenWorkFile(key)
}
