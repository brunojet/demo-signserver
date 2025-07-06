package adapters

import "time"

type HttpMethod string

const (
	HttpMethodGet HttpMethod = "GET"
	HttpMethodPut HttpMethod = "PUT"
)

type StorageServiceInterface interface {
	GetBucketName() string
	GeneratePresignedURL(httpMethod HttpMethod, key string, expires time.Duration) (string, error)
	DownloadFile(key, dest string) error
	UploadFile(key, src string) error
}
