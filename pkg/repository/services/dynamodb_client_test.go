package db_services

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("PROJECT_NAME", "signserver")
	os.Setenv("ENVIRONMENT", "dev")
	os.Setenv("DYNAMODB_ENDPOINT", testResolverEndpoint)
	os.Setenv("AWS_ACCESS_KEY_ID", "fake")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "fake")
}

func getTestDynamoDBClientAndResolver(tableName string) (*dynamodb.Client, *DynamoDBEndpointResolver, error) {
	client, resolver, err := NewDynamoDBClient(context.TODO(), tableName)
	return client, resolver, err
}

func TestNewDynamoDBClient_Success(t *testing.T) {
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8001")
	client, _, err := getTestDynamoDBClientAndResolver("TestTable")
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
