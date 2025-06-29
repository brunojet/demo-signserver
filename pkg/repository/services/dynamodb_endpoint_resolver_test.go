package db_services

import (
	"context"
	"os"
	"testing"
)

var (
	testResolverTableName = "TestTable-002"
	testResolverEndpoint  = "http://localhost:8001"
)

func TestDynamoDBEndpointResolver_EnsureTableExists(t *testing.T) {
	client, resolver, _ := getTestDynamoDBClientAndResolver(t, testResolverTableName)
	ctx := context.TODO()

	// Tenta criar de novo (não deve dar erro)
	err := resolver.EnsureTableExists(ctx, client, "sk")
	if err != nil {
		t.Fatalf("Erro ao garantir tabela já existente: %v", err)
	}
}

func TestDynamoDBEndpointResolver_EnsureTableExists_NoClient(t *testing.T) {
	_, resolver, _ := getTestDynamoDBClientAndResolver(t, testResolverTableName)
	ctx := context.TODO()

	os.Setenv("DELETE_TABLE", "true") // Força deleção da tabela

	// Tenta criar de novo (não deve dar erro)
	err := resolver.EnsureTableExists(ctx, nil, "sk")
	if err == nil {
		t.Fatalf("DynamoDB client erro esperado")
	}
}

func TestDynamoDBEndpointResolver_EnsureTableExists_NoClient_2(t *testing.T) {
	_, resolver, _ := getTestDynamoDBClientAndResolver(t, testResolverTableName)
	ctx := context.TODO()

	os.Setenv("DELETE_TABLE", "") // Força deleção da tabela

	// Tenta criar de novo (não deve dar erro)
	err := resolver.EnsureTableExists(ctx, nil, "sk")
	if err == nil {
		t.Fatalf("DynamoDB client erro esperado")
	}
}
