package db_services

import (
	"testing"
)

var (
	testResolverTableName = "TestTable-002"
	testResolverEndpoint  = "http://localhost:8001"
)

func TestDynamoDBEndpointResolver_EnsureTableExists(t *testing.T) {
	getTestDynamoDBClient(testResolverTableName, "sk")
}
