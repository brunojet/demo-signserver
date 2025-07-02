package db_services

import (
	"context"
	"os"
	"testing"
)

var (
	testTableName = "TestTable-001"
	db            *DBServices
)

func init() {
	os.Setenv("PROJECT_NAME", "signserver")
	os.Setenv("ENVIRONMENT", "dev")
	db = NewDB(testTableName)
}

func TestTableExists_Create_Delete(t *testing.T) {
	table := testTableName

	_ = db.DeleteTable(context.Background(), table)
	if db.TableExists(context.Background(), table) {
		t.Fatalf("Tabela deveria não existir")
	}

	// Cria tabela
	err := db.CreateTable(context.Background(), table, "sk")
	if err != nil {
		t.Fatalf("Erro ao criar tabela: %v", err)
	}
	if !db.TableExists(context.Background(), table) {
		t.Fatalf("Tabela deveria existir após criação")
	}

	// Cria de novo (idempotente)
	err = db.CreateTable(context.Background(), table, "sk")
	if err != nil {
		t.Fatalf("CreateTable deveria ser idempotente: %v", err)
	}

	// Deleta tabela
	err = db.DeleteTable(context.Background(), table)
	if err != nil {
		t.Fatalf("Erro ao deletar tabela: %v", err)
	}
	if db.TableExists(context.Background(), table) {
		t.Fatalf("Tabela deveria não existir")
	}

	// Deleta de novo (idempotente)
	err = db.DeleteTable(context.Background(), table)
	if err != nil {
		t.Fatalf("Erro ao deletar tabela: %v", err)
	}
	if db.TableExists(context.Background(), table) {
		t.Fatalf("Tabela deveria não existir")
	}
}
