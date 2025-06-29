package db_services

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoDBEndpointResolver é um endpoint resolver customizado e utilitário para ambiente local/dev.
type DynamoDBEndpointResolver struct {
	EndpointURL string
	TableName   string
}

// ResolveEndpoint implementa a interface de resolução de endpoint customizado.
func (r *DynamoDBEndpointResolver) ResolveEndpoint(service, region string, options ...interface{}) (aws.Endpoint, error) {
	if service == dynamodb.ServiceID && r.EndpointURL != "" {
		fmt.Printf("[DynamoDBService] Usando endpoint customizado: %s\n", r.EndpointURL)
		return aws.Endpoint{
			URL:           r.EndpointURL,
			SigningRegion: "us-east-1",
		}, nil
	}
	return aws.Endpoint{}, &aws.EndpointNotFoundError{}
}

// EnsureTableExists cria a tabela se ela não existir (útil para dev/test).
// pkKey é obrigatório e será sempre a chave HASH. Se skKey for fornecido, será a RANGE (com nome físico "sk").
func (r *DynamoDBEndpointResolver) EnsureTableExists(ctx context.Context, client *dynamodb.Client, skKey string) error {
	// Se DELETE_TABLE estiver setada, deleta a tabela antes de criar
	if os.Getenv("DELETE_TABLE") != "" {
		err := DeleteTable(ctx, client, r.TableName)
		if err != nil {
			return err
		}
	}

	// Cria tabela
	err := CreateTable(ctx, client, r.TableName, skKey)
	if err != nil {
		return err
	}
	return nil
}
