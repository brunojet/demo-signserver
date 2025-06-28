package storage_services

import (
	"context"
	"errors"
	"testing"
	"time"

	storage_mock "demo-signserver/pkg/storage/mock"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePresignedURL_Success(t *testing.T) {
	svc := &S3Service{Bucket: "bucket"}
	url := "https://example.com/presigned"
	presignClient := &storage_mock.MockPresignClient{URL: url}

	oldNewPresignClient := newPresignClient
	newPresignClient = func(_ *s3.Client) PresignPutObjectAPI { return presignClient }
	defer func() { newPresignClient = oldNewPresignClient }()

	result, err := svc.GeneratePresignedURL("file.apk", 10*time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, url, result)
}

func TestGeneratePresignedURL_Error(t *testing.T) {
	svc := &S3Service{Bucket: "bucket"}
	errMock := errors.New("presign error")
	presignClient := &storage_mock.MockPresignClient{Err: errMock}

	oldNewPresignClient := newPresignClient
	newPresignClient = func(_ *s3.Client) PresignPutObjectAPI { return presignClient }
	defer func() { newPresignClient = oldNewPresignClient }()

	result, err := svc.GeneratePresignedURL("file.apk", 10*time.Minute)
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestNewS3Service_Success(t *testing.T) {
	// Testa inicialização sem erro (mockando config.LoadDefaultConfig)
	svc, err := NewS3Service("bucket")
	assert.NoError(t, err)
	assert.NotNil(t, svc)
	assert.Equal(t, "bucket", svc.Bucket)
}

func TestNewS3Service_Error(t *testing.T) {
	// Simula erro usando injeção de dependência
	mockErr := errors.New("erro de config")
	mockLoader := func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
		return aws.Config{}, mockErr
	}
	_, err := NewS3ServiceWithConfigLoader("bucket", mockLoader)
	assert.Error(t, err)
}

func TestGenerateTemporaryS3Credentials_Success(t *testing.T) {
	mockSTS := &storage_mock.MockSTSClient{
		AssumeRoleFunc: func(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error) {
			return &sts.AssumeRoleOutput{
				Credentials: &types.Credentials{
					AccessKeyId:     aws.String("mock-access-key"),
					SecretAccessKey: aws.String("mock-secret-key"),
					SessionToken:    aws.String("mock-session-token"),
				},
			}, nil
		},
	}
	mockLoader := func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
		return aws.Config{}, nil
	}
	roleArn := "arn:aws:iam::123456789012:role/test-role"
	creds, err := GenerateTemporaryS3CredentialsWithClient(roleArn, "test-session", 900*time.Second, mockLoader, func(cfg aws.Config) STSAPI { return mockSTS })
	assert.NoError(t, err)
	assert.NotNil(t, creds)
	assert.Equal(t, "mock-access-key", *creds.AccessKeyId)
}
