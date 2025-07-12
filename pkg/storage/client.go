package storage

import "time"

type StorageService struct {
	adapter StorageAdapter
}

func NewStorageService(adapter StorageAdapter) StorageAdapter {
	return &StorageService{
		adapter: adapter,
	}
}

func (s *StorageService) DownloadFileFromS3(fileInfo *FileInfo) error {
	return s.adapter.DownloadFileFromS3(fileInfo)
}

func (s *StorageService) GeneratePresignedURL(httpMethod HttpMethod, fileInfo *FileInfo, expires time.Duration) (string, error) {
	return s.adapter.GeneratePresignedURL(httpMethod, fileInfo, expires)
}

func (s *StorageService) GetWorkFilePath(fileInfo *FileInfo) string {
	return s.adapter.GetWorkFilePath(fileInfo)
}

func (s *StorageService) UploadToS3(fileInfo *FileInfo) error {
	return s.adapter.UploadToS3(fileInfo)
}
