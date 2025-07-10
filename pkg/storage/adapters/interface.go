package adapters

import (
	"io"
	"time"
)

type HttpMethod string

const (
	HttpMethodGet HttpMethod = "GET"
	HttpMethodPut HttpMethod = "PUT"
)

type StorageServiceInterface interface {
	GetBucketName() string
	GeneratePresignedURL(httpMethod HttpMethod, key string, expires time.Duration) (string, error)
	DownloadFileFromS3(key string) error
	UploadToS3(key string) error
	OpenWorkFile(key string) (io.ReadWriteCloser, error)
}
