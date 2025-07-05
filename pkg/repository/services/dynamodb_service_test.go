package db_services

import (
	"context"
	"demo-signserver/pkg/repository/domain"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	db_service_test *DynamoDBService
)

func TestMain(m *testing.M) {
	// Setup global (criação de tabela, etc)
	tableName := InitTestTable(NO_KEY)
	defer DeleteTestTable(tableName)

	db_service_test = NewDynamoDBService(tableName, ID_KEY, NO_KEY)

	code := m.Run()
	os.Exit(code)
}

func TestDynamoDBService_CreateItem(t *testing.T) {
	const pk_value = "001"
	type Item struct {
		domain.BaseDomain
	}
	item := Item{}
	item.SetID(pk_value)
	err := db_service_test.CreateItem(context.Background(), &item)
	assert.NoError(t, err)
}

func TestDynamoDBService_CreateItem_Error(t *testing.T) {
	const pk_value = "0011"
	type Item struct {
		domain.BaseDomain
	}
	item := Item{}
	item.SetID(pk_value)
	err := db_service_test.CreateItem(context.Background(), &item)
	assert.NoError(t, err)

	err = db_service_test.CreateItem(context.Background(), &item)
	assert.Error(t, err)

}

func TestDynamoDBService_GetItem(t *testing.T) {
	const pk_value = "123"
	type Item struct {
		domain.BaseDomain
	}
	item := Item{}
	item.SetID(pk_value)
	err := db_service_test.CreateItem(context.Background(), &item)
	assert.NoError(t, err)
	var out Item
	err = db_service_test.GetItem(context.Background(), item.ID, &out)
	assert.NoError(t, err)
	assert.Equal(t, item.ID, out.ID)
}

func TestDynamoDBService_GetItem_Error(t *testing.T) {
	const pk_value = "1234"
	type Item struct {
		domain.BaseDomain
	}
	item := Item{}
	item.SetID(pk_value)
	err := db_service_test.CreateItem(context.Background(), &item)
	assert.NoError(t, err)
	var out Item
	err = db_service_test.GetItem(context.Background(), item.ID+"1", &out)
	assert.Error(t, err)
}

func TestDynamoDBService_UpdateItem(t *testing.T) {
	const pk_value = "456"
	type Item struct {
		domain.BaseDomain
		Name string `dynamodbav:"name"`
	}

	item := Item{Name: "original"}
	item.SetID(pk_value)
	err := db_service_test.CreateItem(context.Background(), &item)
	assert.NoError(t, err)

	// Atualiza campo Name
	update := Item{Name: "updated"}
	err = db_service_test.UpdateItem(context.Background(), item.ID, &update)
	assert.NoError(t, err)

	// Busca e valida
	var out Item
	err = db_service_test.GetItem(context.Background(), item.ID, &out)
	assert.NoError(t, err)
	assert.Equal(t, item.ID, out.ID)
	assert.Equal(t, "updated", out.Name)
}

func TestDynamoDBService_UpdateItem_Error(t *testing.T) {
	const pk_value = "4567"
	type Item struct {
		domain.BaseDomain
		Name string `dynamodbav:"name"`
	}

	item := Item{
		BaseDomain: domain.BaseDomain{ID: pk_value},
		Name:       "original",
	}
	err := db_service_test.CreateItem(context.Background(), &item)
	assert.NoError(t, err)

	// Atualiza campo Name
	update := Item{Name: "updated"}
	err = db_service_test.UpdateItem(context.Background(), item.ID+"1", &update)
	assert.Error(t, err)
}

func TestNewDynamoDBService_Coverage(t *testing.T) {
	tableName := InitTestTable(NO_KEY)
	defer DeleteTestTable(tableName)
	svc := NewDynamoDBService(tableName, ID_KEY, "demo")
	assert.NotNil(t, svc)
	assert.Equal(t, tableName, svc.Table)
	assert.Equal(t, ID_KEY, svc.PKKey)
	assert.Equal(t, "demo", svc.SKKey)
}
