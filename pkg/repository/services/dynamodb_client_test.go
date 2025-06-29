package db_services

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

// Mock para EnsureTableExistsAPI
type mockEnsureTableExists struct {
	fail bool
}

func (m *mockEnsureTableExists) EnsureTableExists(ctx context.Context, client *dynamodb.Client, skKey string) error {
	if m.fail {
		return context.DeadlineExceeded
	}
	return nil
}

func TestNewDynamoDBClient_Success(t *testing.T) {
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8001")
	mock := &mockEnsureTableExists{fail: false}
	client, err := NewDynamoDBClient(context.TODO(), "TestTable", "pk", "sk", mock)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestNewDynamoDBClient_EnsureTableExistsError(t *testing.T) {
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8001")
	mock := &mockEnsureTableExists{fail: true}
	client, err := NewDynamoDBClient(context.TODO(), "TestTable", "pk", "sk", mock)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func mockLoadConfigError(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
	return aws.Config{}, errors.New("mock config error")
}
