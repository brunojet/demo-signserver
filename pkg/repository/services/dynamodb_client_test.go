package db_services

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("PROJECT_NAME", "signserver")
	os.Setenv("ENVIRONMENT", "dev")
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8001")
	os.Setenv("AWS_ACCESS_KEY_ID", "fake")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "fake")
}

func getTestDynamoDBClient(tableName string, sk string) (*dynamodb.Client, error) {
	initTestTable(tableName, sk)
	client, err := NewDynamoDBClient(context.TODO(), tableName)
	return client, err
}

func initTestTable(tableName string, sk string) {
	db := NewDB(tableName)
	err := db.DeleteTable(context.TODO(), tableName)

	if err != nil {
		log.Fatalf("Error deleting table %s: %v", tableName, err)
	}
	err = db.CreateTable(context.TODO(), tableName, sk)

	if err != nil {
		log.Fatalf("Error creating table %s: %v", tableName, err)
	}
}

func TestNewDynamoDBClient_Success(t *testing.T) {
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8001")
	client, err := getTestDynamoDBClient("TestTable", "")
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
