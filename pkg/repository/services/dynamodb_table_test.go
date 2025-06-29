package db_services

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func init() {
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8001")
	os.Setenv("AWS_ACCESS_KEY_ID", "fake")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "fake")
}

func getTestDynamoDBClient(t *testing.T) *dynamodb.Client {
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		t.Skip("DYNAMODB_ENDPOINT não setado para testes locais")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolver(
			customTestResolver(endpoint),
		),
	)
	if err != nil {
		t.Fatalf("Erro ao carregar config: %v", err)
	}
	return dynamodb.NewFromConfig(cfg)
}

type customTestResolver string

func (r customTestResolver) ResolveEndpoint(service, region string) (aws.Endpoint, error) {
	return aws.Endpoint{URL: string(r), SigningRegion: "us-east-1"}, nil
}

func TestTableExists_Create_Delete(t *testing.T) {
	table := "test-table-dynamodb-table"
	client := getTestDynamoDBClient(t)
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
