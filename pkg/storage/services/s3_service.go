package storage_services

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Service encapsula operações com o S3.
type S3Service struct {
	Client *s3.Client
	Bucket string
}

// NewS3ServiceWithConfigLoader permite injetar função de carregamento de config (para testes).
func NewS3ServiceWithConfigLoader(bucket string, loadConfig func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error)) (*S3Service, error) {
	cfg, err := loadConfig(context.Background())
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	return &S3Service{Client: client, Bucket: bucket}, nil
}

// NewS3Service padrão, usa config.LoadDefaultConfig
func NewS3Service(bucket string) *S3Service {
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
func (s *S3Service) GeneratePresignedPutURL(key string, expires time.Duration) (string, error) {
	presignClient := newPresignClient(s.Client)
	params := &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}
	presignedReq, err := presignClient.PresignPutObject(context.Background(), params, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})
	if err != nil {
		return "", err
	}
	return presignedReq.URL, nil
}

// Gera uma URL pré-assinada para GET (download) de um objeto S3
func (s *S3Service) GeneratePresignedGetURL(key string, expires time.Duration) (string, error) {
	presignClient := newPresignClient(s.Client)
	params := &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}
	presignedReq, err := presignClient.PresignGetObject(context.Background(), params, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})
	if err != nil {
		return "", err
	}
	return presignedReq.URL, nil
}

// S3ServiceInterface define as operações expostas para uso/mocks
// e abstrai detalhes de implementação do S3 real.
type S3ServiceInterface interface {
	GeneratePresignedURL(key string) (string, error)
	GeneratePresignedURLWithExpiry(key string, expires time.Duration) (string, error)
}
