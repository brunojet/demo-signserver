package storages

import (
	"demo-signserver/internal/config"
	"demo-signserver/pkg/storage"
	"time"
)

type StorageService struct {
	storage storage.StorageAdapter
}

func NewStorageService() storage.StorageAdapter {
	methods := config.GetSignServerMethods()
	return &StorageService{
		storage: methods.NewStorageService(),
	}
}

func (s *StorageService) GeneratePresignedURL(httpMethod storage.HttpMethod, fileInfo *storage.FileInfo, expires time.Duration) (string, error) {
	return s.storage.GeneratePresignedURL(httpMethod, fileInfo, expires)
}

// DownloadFileFromS3 implements adapters.StorageServiceInterface.
func (s *StorageService) DownloadFileFromS3(fileInfo *storage.FileInfo) error {
	return s.storage.DownloadFileFromS3(fileInfo)
}

// UploadToS3 implements adapters.StorageServiceInterface.
func (s *StorageService) UploadToS3(fileInfo *storage.FileInfo) error {
	return s.storage.UploadToS3(fileInfo)
}

func (s *StorageService) GetWorkFilePath(fileInfo *storage.FileInfo) string {
	return s.storage.GetWorkFilePath(fileInfo)
}
