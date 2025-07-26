package adapters

import (
	"context"
	"demo-signserver/pkg/file"
	"demo-signserver/pkg/storage"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var _ storage.StorageAdapter = (*S3Service)(nil)

type S3Service struct {
	Client   *s3.Client
	Bucket   string
	workPath string
}

func (s *S3Service) DownloadFileFromS3(fileInfo *storage.FileInfo) error {
	input := &s3.GetObjectInput{
		Bucket: aws.String(fileInfo.StoragePath),
		Key:    aws.String(fileInfo.FilePath),
	}
	resp, err := s.Client.GetObject(context.Background(), input)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dst := s.GetWorkFilePath(fileInfo)

	f, err := file.CreateFilePath(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

func (s *S3Service) GeneratePresignedURL(httpMethod storage.HttpMethod, fileInfo *storage.FileInfo, expires time.Duration) (string, error) {
	var err error
	var presignedResponse *v4.PresignedHTTPRequest
	switch httpMethod {
	case storage.HttpMethodPut:
		presignedResponse, err = s.generatePresignedPutURL(fileInfo.FilePath, expires)
	case storage.HttpMethodGet:
		presignedResponse, err = s.generatePresignedGetURL(fileInfo.FilePath, expires)
	}
	if err != nil {
		return "", fmt.Errorf("erro ao gerar URL pré-assinada: %w", err)
	}
	return presignedResponse.URL, nil
}

func (s *S3Service) GetWorkFilePath(fileInfo *storage.FileInfo) string {
	return filepath.Join(s.workPath, fileInfo.FilePath)
}

func (s *S3Service) UploadToS3(fileInfo *storage.FileInfo) error {
	src := s.GetWorkFilePath(fileInfo)
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	input := &s3.PutObjectInput{
		Bucket: aws.String(fileInfo.StoragePath),
		Key:    aws.String(fileInfo.FilePath),
		Body:   f,
	}
	_, err = s.Client.PutObject(context.Background(), input)
	if err != nil {
		log.Printf("[S3Service] Erro ao fazer upload do arquivo %s para o bucket %s: %v\n", fileInfo.FilePath, fileInfo.StoragePath, err)
		return err
	}

	fileInfo.State = storage.FileStateReady
	fileInfo.Size, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		log.Printf("[S3Service] Erro ao obter tamanho do arquivo %s: %v\n", fileInfo.FilePath, err)
	}

	return nil
}

func NewS3ServiceWithConfigLoader(bucket string, loadConfig func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error)) (storage.StorageAdapter, error) {
	workPath := filepath.Join(os.TempDir(), "s3_storage_service")
	file.CreatePathIfNotExists(workPath)

	cfg, err := loadConfig(context.Background())
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	return &S3Service{Client: client, Bucket: bucket, workPath: workPath}, nil
}

// NewS3Service padrão, usa config.LoadDefaultConfig
func NewS3Service(bucket string) storage.StorageAdapter {
	s3Service, err := NewS3ServiceWithConfigLoader(bucket, config.LoadDefaultConfig)
	if err != nil {
		log.Fatalf("[S3Service] Erro ao criar S3Service: %v", err)
	}
	return s3Service
}

var newPresignClient = func(client *s3.Client) PresignObjectAPI {
	return s3.NewPresignClient(client)
}

// PresignObjectAPI define a interface para mocks do PresignClient do S3.
type PresignObjectAPI interface {
	PresignPutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
	PresignGetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

// Ajusta GeneratePresignedURL para usar a interface e permitir mock nos testes
func (s *S3Service) generatePresignedPutURL(key string, expires time.Duration) (*v4.PresignedHTTPRequest, error) {
	presignClient := newPresignClient(s.Client)
	params := &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}
	return presignClient.PresignPutObject(context.Background(), params, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})
}

// Gera uma URL pré-assinada para GET (download) de um objeto S3
func (s *S3Service) generatePresignedGetURL(key string, expires time.Duration) (*v4.PresignedHTTPRequest, error) {
	presignClient := newPresignClient(s.Client)
	params := &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}
	return presignClient.PresignGetObject(context.Background(), params, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})
}

// S3ServiceInterface define as operações expostas para uso/mocks
// e abstrai detalhes de implementação do S3 real.
type S3ServiceInterface interface {
	GeneratePresignedURL(key string) (string, error)
	GeneratePresignedURLWithExpiry(key string, expires time.Duration) (string, error)
}
