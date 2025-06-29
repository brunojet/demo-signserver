package db_services

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

var (
	dynamoClient *dynamodb.Client
	tableName    = "TestTable"
	endpoint     = "http://localhost:8001"
)

func init() {
	os.Setenv("DYNAMODB_ENDPOINT", endpoint)
	os.Setenv("AWS_ACCESS_KEY_ID", "fake")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "fake")
	os.Setenv("DELETE_TABLE", "true")
}

func TestMain(m *testing.M) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolver(
			customTestResolver(endpoint),
		),
	)
	if err != nil {
		panic(err)
	}
	dynamoClient = dynamodb.NewFromConfig(cfg)

	// Cria a tabela antes dos testes
	err = CreateTable(context.TODO(), dynamoClient, tableName, "")
	if err != nil {
		panic(err)
	}

	code := m.Run()

	// Deleta a tabela após os testes
	_ = DeleteTable(context.TODO(), dynamoClient, tableName)
	os.Exit(code)
}

func TestDynamoDBService_CreateItem(t *testing.T) {
	const pk_value = "001"
	type Item struct {
		Name string `dynamodbav:"name"`
	}
	service := &DynamoDBService{Client: dynamoClient, Table: tableName, PKKey: "name", SKKey: ""}
	item := Item{Name: pk_value}
	err := service.CreateItem(context.TODO(), item)
	assert.NoError(t, err)
}

func TestDynamoDBService_GetItem(t *testing.T) {
	const pk_value = "123"
	type Item struct {
		PK string `dynamodbav:"pk"`
	}
	service := &DynamoDBService{Client: dynamoClient, Table: tableName, PKKey: "pk", SKKey: ""}
	item := Item{PK: pk_value}
	_ = service.CreateItem(context.TODO(), item)
	var out Item
	err := service.GetItem(context.TODO(), pk_value, &out)
	assert.NoError(t, err)
	assert.Equal(t, pk_value, out.PK)
}

func TestDynamoDBService_UpdateItem(t *testing.T) {
	const pk_value = "456"
	type Item struct {
		PK   string `dynamodbav:"pk"`
		Name string `dynamodbav:"name"`
	}

	service := &DynamoDBService{Client: dynamoClient, Table: tableName, PKKey: "pk", SKKey: ""}
	// Cria item inicial
	item := Item{PK: pk_value, Name: "original"}
	err := service.CreateItem(context.TODO(), item)
	assert.NoError(t, err)

	// Atualiza campo Name
	update := Item{Name: "updated"}
	err = service.UpdateItem(context.TODO(), pk_value, update)
	assert.NoError(t, err)

	// Busca e valida
	var out Item
	err = service.GetItem(context.TODO(), pk_value, &out)
	assert.NoError(t, err)
	assert.Equal(t, pk_value, out.PK)
	assert.Equal(t, "updated", out.Name)
}
