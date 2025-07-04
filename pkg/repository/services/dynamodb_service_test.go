package db_services

import (
	"context"
	"demo-signserver/pkg/repository/domain"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	service          *DynamoDBService
	tableNameService = "TableService"
)

func init() {
	initTestTable(tableNameService, "")
	var err error
	service, err = NewDynamoDBService(tableNameService, "pk", "")
	if err != nil {
		log.Fatalf("Error initializing DynamoDB service: %v\n", err)
	}
}

func TestDynamoDBService_CreateItem(t *testing.T) {
	const pk_value = "001"
	type Item struct {
		domain.BaseDomain
		Name string `dynamodbav:"name"`
	}
	service, err := NewDynamoDBService(tableNameService, "name", "")
	assert.NoError(t, err)
	item := Item{Name: pk_value}
	err = service.CreateItem(context.Background(), &item)
	assert.NoError(t, err)
}

func TestDynamoDBService_CreateItem_Error(t *testing.T) {
	const pk_value = "0011"
	type Item struct {
		domain.BaseDomain
		Name string `dynamodbav:"name"`
	}
	service, err := NewDynamoDBService(tableNameService, "name", "")
	assert.NoError(t, err)
	item := Item{Name: pk_value}
	err = service.CreateItem(context.Background(), &item)
	assert.NoError(t, err)

	err = service.CreateItem(context.Background(), &item)
	assert.Error(t, err)

}

func TestDynamoDBService_GetItem(t *testing.T) {
	const pk_value = "123"
	type Item struct {
		domain.BaseDomain
		PK string `dynamodbav:"pk"`
	}

	item := Item{PK: pk_value}
	err := service.CreateItem(context.Background(), &item)
	assert.NoError(t, err)
	var out Item
	err = service.GetItem(context.Background(), item.ID, &out)
	assert.NoError(t, err)
	assert.Equal(t, item.ID, out.PK)
}

func TestDynamoDBService_GetItem_Error(t *testing.T) {
	const pk_value = "1234"
	type Item struct {
		domain.BaseDomain
		PK string `dynamodbav:"pk"`
	}

	item := Item{PK: pk_value}
	err := service.CreateItem(context.Background(), &item)
	assert.NoError(t, err)
	var out Item
	err = service.GetItem(context.Background(), item.ID+"1", &out)
	assert.Error(t, err)
}

func TestDynamoDBService_UpdateItem(t *testing.T) {
	const pk_value = "456"
	type Item struct {
		domain.BaseDomain
		PK   string `dynamodbav:"pk"`
		Name string `dynamodbav:"name"`
	}

	item := Item{PK: pk_value, Name: "original"}
	err := service.CreateItem(context.Background(), &item)
	assert.NoError(t, err)

	// Atualiza campo Name
	update := Item{Name: "updated"}
	err = service.UpdateItem(context.Background(), item.ID, &update)
	assert.NoError(t, err)

	// Busca e valida
	var out Item
	err = service.GetItem(context.Background(), item.ID, &out)
	assert.NoError(t, err)
	assert.Equal(t, item.ID, out.PK)
	assert.Equal(t, "updated", out.Name)
}

func TestDynamoDBService_UpdateItem_Error(t *testing.T) {
	const pk_value = "4567"
	type Item struct {
		domain.BaseDomain
		PK   string `dynamodbav:"pk"`
		Name string `dynamodbav:"name"`
	}

	item := Item{PK: pk_value, Name: "original"}
	err := service.CreateItem(context.Background(), &item)
	assert.NoError(t, err)

	// Atualiza campo Name
	update := Item{Name: "updated"}
	err = service.UpdateItem(context.Background(), item.ID+"1", &update)
	assert.Error(t, err)
}

func TestNewDynamoDBService_Coverage(t *testing.T) {
	svc, err := NewDynamoDBService(tableNameService, "pk", "sk")
	assert.NoError(t, err)
	assert.NotNil(t, svc)
	assert.Equal(t, fmt.Sprintf("%s-%s-%s", os.Getenv("PROJECT_NAME"), os.Getenv("ENVIRONMENT"), tableNameService), svc.Table)
	assert.Equal(t, "pk", svc.PKKey)
	assert.Equal(t, "sk", svc.SKKey)
}
