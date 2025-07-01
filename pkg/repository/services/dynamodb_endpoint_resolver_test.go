package db_services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testResolverTableName = "TestTable-002"
	testResolverEndpoint  = "http://localhost:8001"
)

func TestDynamoDBEndpointResolver_EnsureTableExists(t *testing.T) {
	client, err := getTestDynamoDBClient(testResolverTableName, "sk")
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
