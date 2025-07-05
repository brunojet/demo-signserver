package db_services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableExists_Create_Delete(t *testing.T) {
	table := "TestTable-001"

	_ = table_service_test.DeleteTable(context.Background(), table)
	tableExists := table_service_test.TableExists(context.Background(), table)
	assert.False(t, tableExists, "Tabela deveria não existir antes de criar")

	// Cria tabela
	err := table_service_test.CreateTable(context.Background(), table, SORT_KEY)
	assert.NoError(t, err, "Erro ao criar tabela")

	tableExists = table_service_test.TableExists(context.Background(), table)
	assert.True(t, tableExists, "Tabela deveria existir após criação")

	// Cria de novo (idempotente)
	err = table_service_test.CreateTable(context.Background(), table, SORT_KEY)
	assert.NoError(t, err, "CreateTable deveria ser idempotente")

	// Deleta tabela
	err = table_service_test.DeleteTable(context.Background(), table)
	assert.NoError(t, err, "Erro ao deletar tabela")

	tableExists = table_service_test.TableExists(context.Background(), table)
	assert.False(t, tableExists, "Tabela deveria não existir após deleção")

	// Deleta de novo (idempotente)
	err = table_service_test.DeleteTable(context.Background(), table)
	assert.NoError(t, err, "DeleteTable deveria ser idempotente")

	tableExists = table_service_test.TableExists(context.Background(), table)
	assert.False(t, tableExists, "Tabela deveria não existir após deleção")
}
