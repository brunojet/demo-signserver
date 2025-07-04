package storage_services

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

// STSAPI define interface para mocks do STS Client.
type STSAPI interface {
	AssumeRole(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error)
}

// Função para criar o STS client, permite injeção de mock em testes
var newSTSClient = func(cfg aws.Config) STSAPI {
	return sts.NewFromConfig(cfg)
}

// Gera credenciais temporárias STS para escrita no S3
func GenerateTemporaryS3CredentialsWithClient(roleArn, sessionName string, duration time.Duration, loadConfig func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error), stsClientFactory func(cfg aws.Config) STSAPI) (*types.Credentials, error) {
	cfg, err := loadConfig(context.Background())
	if err != nil {
		return nil, err
	}
	stsClient := stsClientFactory(cfg)
	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(sessionName),
		DurationSeconds: aws.Int32(int32(duration.Seconds())),
	}
	result, err := stsClient.AssumeRole(context.Background(), input)
	if err != nil {
		return nil, err
	}
	return result.Credentials, nil
}

// Versão padrão, usa config.LoadDefaultConfig e newSTSClient
func GenerateTemporaryS3Credentials(roleArn, sessionName string, duration time.Duration) (*types.Credentials, error) {
	return GenerateTemporaryS3CredentialsWithClient(roleArn, sessionName, duration, config.LoadDefaultConfig, newSTSClient)
}
