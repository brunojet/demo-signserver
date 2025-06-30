package storage_services

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

// S3Service encapsula operações com o S3.
type S3Service struct {
	Client *s3.Client
	Bucket string
}

// NewS3ServiceWithConfigLoader permite injetar função de carregamento de config (para testes).
func NewS3ServiceWithConfigLoader(bucket string, loadConfig func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error)) (*S3Service, error) {
	cfg, err := loadConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	return &S3Service{Client: client, Bucket: bucket}, nil
}

// NewS3Service padrão, usa config.LoadDefaultConfig
func NewS3Service(bucket string) (*S3Service, error) {
	return NewS3ServiceWithConfigLoader(bucket, config.LoadDefaultConfig)
}

var newPresignClient = func(client *s3.Client) PresignPutObjectAPI {
	return s3.NewPresignClient(client)
}

// PresignPutObjectAPI define a interface para mocks do PresignClient do S3.
type PresignPutObjectAPI interface {
	PresignPutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

// Ajusta GeneratePresignedURL para usar a interface e permitir mock nos testes
func (s *S3Service) GeneratePresignedURL(key string, expires time.Duration) (string, error) {
	presignClient := newPresignClient(s.Client)
	params := &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}
	presignedReq, err := presignClient.PresignPutObject(context.TODO(), params, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})
	if err != nil {
		return "", err
	}
	return presignedReq.URL, nil
}

// Função para criar o STS client, permite injeção de mock em testes
var newSTSClient = func(cfg aws.Config) STSAPI {
	return sts.NewFromConfig(cfg)
}

// STSAPI define interface para mocks do STS Client.
type STSAPI interface {
	AssumeRole(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error)
}

// Gera credenciais temporárias STS para escrita no S3
func GenerateTemporaryS3CredentialsWithClient(roleArn, sessionName string, duration time.Duration, loadConfig func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error), stsClientFactory func(cfg aws.Config) STSAPI) (*types.Credentials, error) {
	cfg, err := loadConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	stsClient := stsClientFactory(cfg)
	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(sessionName),
		DurationSeconds: aws.Int32(int32(duration.Seconds())),
	}
	result, err := stsClient.AssumeRole(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return result.Credentials, nil
}

// Versão padrão, usa config.LoadDefaultConfig e newSTSClient
func GenerateTemporaryS3Credentials(roleArn, sessionName string, duration time.Duration) (*types.Credentials, error) {
	return GenerateTemporaryS3CredentialsWithClient(roleArn, sessionName, duration, config.LoadDefaultConfig, newSTSClient)
}
