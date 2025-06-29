package db_services

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	testResolverTableName = "TestTable-002"
	testResolverEndpoint  = "http://localhost:8001"
)

func init() {
	os.Setenv("DYNAMODB_ENDPOINT", testResolverEndpoint)
	os.Setenv("AWS_ACCESS_KEY_ID", "fake")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "fake")
}

func getTestDynamoDBClientAndResolver(t *testing.T, tableName string) (*dynamodb.Client, *DynamoDBEndpointResolver) {
	resolver := &DynamoDBEndpointResolver{
		EndpointURL: testResolverEndpoint,
		TableName:   tableName,
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.EndpointResolverWithOptions = resolver
		return nil
	})
	if err != nil {
		t.Fatalf("Erro ao carregar config: %v", err)
	}
	client := dynamodb.NewFromConfig(cfg)
	return client, resolver
}

func TestDynamoDBEndpointResolver_EnsureTableExists(t *testing.T) {
	client, resolver := getTestDynamoDBClientAndResolver(t, testResolverTableName)
	ctx := context.TODO()

	// Tenta criar de novo (não deve dar erro)
	err := resolver.EnsureTableExists(ctx, client, "sk")
	if err != nil {
		t.Fatalf("Erro ao garantir tabela já existente: %v", err)
	}
}

func TestDynamoDBEndpointResolver_EnsureTableExists_NoClient(t *testing.T) {
	_, resolver := getTestDynamoDBClientAndResolver(t, testResolverTableName)
	ctx := context.TODO()

	os.Setenv("DELETE_TABLE", "true") // Força deleção da tabela

	// Tenta criar de novo (não deve dar erro)
	err := resolver.EnsureTableExists(ctx, nil, "sk")
	if err == nil {
		t.Fatalf("DynamoDB client erro esperado")
	}
}

func TestDynamoDBEndpointResolver_EnsureTableExists_NoClient_2(t *testing.T) {
	_, resolver := getTestDynamoDBClientAndResolver(t, testResolverTableName)
	ctx := context.TODO()

	os.Setenv("DELETE_TABLE", "") // Força deleção da tabela

	// Tenta criar de novo (não deve dar erro)
	err := resolver.EnsureTableExists(ctx, nil, "sk")
	if err == nil {
		t.Fatalf("DynamoDB client erro esperado")
	}
}
