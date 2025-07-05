package db_services

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8001")
	os.Setenv("AWS_ACCESS_KEY_ID", "fake")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "fake")
}

func CreateTestTable(tableName string, sk string) {
	db := NewDB()
	err := db.CreateTable(context.Background(), tableName, sk)
	if err != nil {
		log.Fatalf("Error creating table %s: %v", tableName, err)
	}
}

func DeleteTestTable(tableName string) {
	db := NewDB()
	err := db.DeleteTable(context.Background(), tableName)
	if err != nil {
		log.Fatalf("Error deleting table %s: %v", tableName, err)
	}
}

func InitTestTable(sk string) string {
	tableName := uuid.New().String()
	CreateTestTable(tableName, sk)
	return tableName
}

func getTestDynamoDBClient(sk string) (*dynamodb.Client, string) {
	tableName := InitTestTable(sk)
	client := GetDynamoDBCLient()
	return client, tableName
}

func TestNewDynamoDBClient_Success(t *testing.T) {
	client, tableName := getTestDynamoDBClient(NO_KEY)
	defer DeleteTestTable(tableName)
	assert.NotNil(t, client)
}
