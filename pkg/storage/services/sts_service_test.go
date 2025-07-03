package storage_services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTemporaryS3Credentials_Success(t *testing.T) {
	mockSTS := &MockSTSClient{
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

func TestGenerateTemporaryS3Credentials_Error(t *testing.T) {
	mockSTS := &MockSTSClient{
		AssumeRoleFunc: func(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error) {
			return nil, errors.New("assume role error")
		},
	}
	mockLoader := func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
		return aws.Config{}, nil
	}
	roleArn := "arn:aws:iam::123456789012:role/test-role"
	creds, err := GenerateTemporaryS3CredentialsWithClient(roleArn, "test-session", 900*time.Second, mockLoader, func(cfg aws.Config) STSAPI { return mockSTS })
	assert.Error(t, err)
	assert.Nil(t, creds)
}
