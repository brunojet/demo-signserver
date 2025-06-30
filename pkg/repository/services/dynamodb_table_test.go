package db_services

import (
	"context"
	"testing"
)

var (
	testTableName = "TestTable-001"
)

func TestTableExists_Create_Delete(t *testing.T) {
	table := testTableName
	client, _, _ := getTestDynamoDBClientAndResolver(testTableName)
	ctx := context.TODO()

	// Garante que a tabela não existe
	_ = DeleteTable(ctx, client, table)
	if TableExists(ctx, client, table) {
		t.Fatalf("Tabela deveria não existir")
	}

	// Cria tabela
	err := CreateTable(ctx, client, table, "sk")
	if err != nil {
		t.Fatalf("Erro ao criar tabela: %v", err)
	}
	if !TableExists(ctx, client, table) {
		t.Fatalf("Tabela deveria existir após criação")
	}

	// Cria de novo (idempotente)
	err = CreateTable(ctx, client, table, "sk")
	if err != nil {
		t.Fatalf("CreateTable deveria ser idempotente: %v", err)
	}

	// Deleta tabela
	err = DeleteTable(ctx, client, table)
	if err != nil {
		t.Fatalf("Erro ao deletar tabela: %v", err)
	}
	if TableExists(ctx, client, table) {
		t.Fatalf("Tabela deveria não existir após deleção")
	}

	// Deleta de novo (idempotente)
	err = DeleteTable(ctx, client, table)
	if err != nil {
		t.Fatalf("DeleteTable deveria ser idempotente: %v", err)
	}
}
