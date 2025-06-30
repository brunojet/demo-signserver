package db_services

import (
	"context"
	"demo-signserver/pkg/repository/domain"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

var (
	dynamoClient     *dynamodb.Client
	tableNameService = "TableService"
)

func init() {
	os.Setenv("DELETE_TABLE", "true")
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8001")
}

func TestMain(m *testing.M) {
	// Use getTestDynamoDBClientAndResolver directly since *testing.T is not available in TestMain
	dynamoClient, _, _ = getTestDynamoDBClientAndResolver(tableNameService)

	// Cria a tabela antes dos testes
	err := CreateTable(context.TODO(), dynamoClient, tableNameService, "")
	if err != nil {
		panic(err)
	}

	code := m.Run()

	// Deleta a tabela após os testes
	_ = DeleteTable(context.TODO(), dynamoClient, tableNameService)
	os.Exit(code)
}

func TestDynamoDBService_CreateItem(t *testing.T) {
	const pk_value = "001"
	type Item struct {
		Name string `dynamodbav:"name"`
	}
	service := &DynamoDBService{Client: dynamoClient, Table: tableNameService, PKKey: "name", SKKey: ""}
	item := Item{Name: pk_value}
	_, err := service.CreateItem(context.TODO(), item)
	assert.NoError(t, err)
}

func TestDynamoDBService_CreateItem_Error(t *testing.T) {
	const pk_value = "0011"
	type Item struct {
		Name string `dynamodbav:"name"`
	}
	service := &DynamoDBService{Client: dynamoClient, Table: tableNameService, PKKey: "name", SKKey: ""}
	item := Item{Name: pk_value}
	_, err := service.CreateItem(context.TODO(), item)
	assert.NoError(t, err)

	_, err = service.CreateItem(context.TODO(), item)
	assert.Error(t, err)

}

func TestDynamoDBService_GetItem(t *testing.T) {
	const pk_value = "123"
	type Item struct {
		domain.BaseDomain
		PK string `dynamodbav:"pk"`
	}
	service := &DynamoDBService{Client: dynamoClient, Table: tableNameService, PKKey: "pk", SKKey: ""}
	item := Item{PK: pk_value}
	pk, _ := service.CreateItem(context.TODO(), item)
	var out Item
	err := service.GetItem(context.TODO(), pk, &out)
	assert.NoError(t, err)
	assert.Equal(t, pk, out.PK)
}

func TestDynamoDBService_GetItem_Error(t *testing.T) {
	const pk_value = "1234"
	type Item struct {
		domain.BaseDomain
		PK string `dynamodbav:"pk"`
	}
	service := &DynamoDBService{Client: dynamoClient, Table: tableNameService, PKKey: "pk", SKKey: ""}
	item := Item{PK: pk_value}
	pk, _ := service.CreateItem(context.TODO(), item)
	var out Item
	err := service.GetItem(context.TODO(), pk+"1", &out)
	assert.Error(t, err)
}

func TestDynamoDBService_UpdateItem(t *testing.T) {
	const pk_value = "456"
	type Item struct {
		domain.BaseDomain
		PK   string `dynamodbav:"pk"`
		Name string `dynamodbav:"name"`
	}

	service := &DynamoDBService{Client: dynamoClient, Table: tableNameService, PKKey: "pk", SKKey: ""}
	// Cria item inicial
	item := Item{PK: pk_value, Name: "original"}
	pk, err := service.CreateItem(context.TODO(), item)
	assert.NoError(t, err)

	// Atualiza campo Name
	update := Item{Name: "updated"}
	err = service.UpdateItem(context.TODO(), pk, update)
	assert.NoError(t, err)

	// Busca e valida
	var out Item
	err = service.GetItem(context.TODO(), pk, &out)
	assert.NoError(t, err)
	assert.Equal(t, pk, out.PK)
	assert.Equal(t, "updated", out.Name)
}

func TestDynamoDBService_UpdateItem_Error(t *testing.T) {
	const pk_value = "4567"
	type Item struct {
		PK   string `dynamodbav:"pk"`
		Name string `dynamodbav:"name"`
	}

	service := &DynamoDBService{Client: dynamoClient, Table: tableNameService, PKKey: "pk", SKKey: ""}
	// Cria item inicial
	item := Item{PK: pk_value, Name: "original"}
	pk, err := service.CreateItem(context.TODO(), item)
	assert.NoError(t, err)

	// Atualiza campo Name
	update := Item{Name: "updated"}
	err = service.UpdateItem(context.TODO(), pk+"1", update)
	assert.Error(t, err)
}

func TestNewDynamoDBService_Coverage(t *testing.T) {
	svc, err := NewDynamoDBService(tableNameService, "pk", "sk")
	assert.NoError(t, err)
	assert.NotNil(t, svc)
	assert.Equal(t, tableNameService, svc.Table)
	assert.Equal(t, "pk", svc.PKKey)
	assert.Equal(t, "sk", svc.SKKey)
}
