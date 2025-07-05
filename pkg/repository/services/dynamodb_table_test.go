package db_services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	table_test_db *DBServices
)

func init() {
	table_test_db = NewDB()
}

func TestTableExists_Create_Delete(t *testing.T) {
	table := "TestTable-001"

	_ = table_test_db.DeleteTable(context.Background(), table)
	tableExists := table_test_db.TableExists(context.Background(), table)
	assert.False(t, tableExists, "Tabela deveria não existir antes de criar")

	// Cria tabela
	err := table_test_db.CreateTable(context.Background(), table, SORT_KEY)
	assert.NoError(t, err, "Erro ao criar tabela")

	tableExists = table_test_db.TableExists(context.Background(), table)
	assert.True(t, tableExists, "Tabela deveria existir após criação")

	// Cria de novo (idempotente)
	err = table_test_db.CreateTable(context.Background(), table, SORT_KEY)
	assert.NoError(t, err, "CreateTable deveria ser idempotente")

	// Deleta tabela
	err = table_test_db.DeleteTable(context.Background(), table)
	assert.NoError(t, err, "Erro ao deletar tabela")

	tableExists = table_test_db.TableExists(context.Background(), table)
	assert.False(t, tableExists, "Tabela deveria não existir após deleção")

	// Deleta de novo (idempotente)
	err = table_test_db.DeleteTable(context.Background(), table)
	assert.NoError(t, err, "DeleteTable deveria ser idempotente")

	tableExists = table_test_db.TableExists(context.Background(), table)
	assert.False(t, tableExists, "Tabela deveria não existir após deleção")
}
