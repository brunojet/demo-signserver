package db_services

import (
	"context"
	"os"
	"testing"
)

var (
	dynamo_service_test *DynamoDBService
	table_service_test  *TableServices
)

func initConfig() {
	os.Setenv("PROJECT_NAME", "signserver")
	os.Setenv("ENVIRONMENT", "dev")
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8001")
	os.Setenv("AWS_ACCESS_KEY_ID", "fake")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "fake")
}

func TestMain(m *testing.M) {
	initConfig()
	tableDbTest := "tableDbTest"
	table_service_test = NewDB()
	code := 1
	defer func() {
		table_service_test.DeleteTable(context.Background(), tableDbTest)
		os.Exit(code)
	}()
	table_service_test.DeleteTable(context.Background(), tableDbTest)
	table_service_test.CreateTable(context.Background(), tableDbTest, NO_KEY)
	dynamo_service_test = NewDynamoDBService(tableDbTest, ID_KEY, NO_KEY)
	code = m.Run()
}
