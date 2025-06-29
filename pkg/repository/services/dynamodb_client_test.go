package db_services

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDynamoDBClient_Success(t *testing.T) {
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8001")
	client, _, err := NewDynamoDBClient(context.TODO(), "TestTable")
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
