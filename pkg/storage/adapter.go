package storage

import (
	"time"
)

type HttpMethod string

type FileState string

const (
	HttpMethodGet HttpMethod = "GET"
	HttpMethodPut HttpMethod = "PUT"
)

const (
	FileStatePending FileState = "pending"
	FileStateReady   FileState = "ready"
)

type FileInfo struct {
	StoragePath string    `dynamodbav:"storage_path" required:"true"`
	FilePath    string    `dynamodbav:"file_path" required:"true"`
	Size        int64     `dynamodbav:"size" default:"0"`
	Hash        string    `dynamodbav:"hash" default:""`
	State       FileState `dynamodbav:"state" default:"pending"`
}

type StorageAdapter interface {
	GetWorkFilePath(fileInfo *FileInfo) string
	GeneratePresignedURL(httpMethod HttpMethod, fileInfo *FileInfo, expires time.Duration) (string, error)
	DownloadFileFromS3(fileInfo *FileInfo) error
	UploadToS3(fileInfo *FileInfo) error
}
