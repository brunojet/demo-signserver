package db_services

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func getTestDynamoDBClientAndResolver(t *testing.T, table string) (*dynamodb.Client, *DynamoDBEndpointResolver) {
	endpoint := "http://localhost:8001"
	if endpoint == "" {
		t.Skip("DYNAMODB_ENDPOINT não setado para testes locais")
	}
	resolver := &DynamoDBEndpointResolver{
		EndpointURL: endpoint,
		TableName:   table,
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
	table := "test-table-dynamodb-endpoint-resolver"
	client, resolver := getTestDynamoDBClientAndResolver(t, table)
	ctx := context.TODO()

	// Tenta criar de novo (não deve dar erro)
	err := resolver.EnsureTableExists(ctx, client, "sk")
	if err != nil {
		t.Fatalf("Erro ao garantir tabela já existente: %v", err)
	}
}
