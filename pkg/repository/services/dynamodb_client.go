package db_services

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// Interface para permitir mock de EnsureTableExists
// (ou crie manualmente)
//
//go:generate mockgen -destination=dynamodb_client_mock.go -package=db_services . EnsureTableExistsAPI
type EnsureTableExistsAPI interface {
	EnsureTableExists(ctx context.Context, client *dynamodb.Client, skKey string) error
}

// NewDynamoDBClient agora aceita injeção do resolver (default: DynamoDBEndpointResolver)
func NewDynamoDBClient(ctx context.Context, table string, pkKey string, skKey string, resolverOpt ...EnsureTableExistsAPI) (*dynamodb.Client, error) {
	loadConfig := config.LoadDefaultConfig
	var resolver EnsureTableExistsAPI
	var realResolver *DynamoDBEndpointResolver
	if len(resolverOpt) > 0 {
		resolver = resolverOpt[0]
	} else {
		realResolver = &DynamoDBEndpointResolver{}
		resolver = realResolver
	}
	cfg, err := loadConfig(ctx, func(o *config.LoadOptions) error {
		if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
			if realResolver == nil {
				realResolver = &DynamoDBEndpointResolver{}
			}
			realResolver.EndpointURL = endpoint
			realResolver.TableName = table
			o.EndpointResolverWithOptions = realResolver
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg)
	// Se for endpoint customizado, garante que a tabela existe
	if realResolver != nil && realResolver.EndpointURL != "" {
		err = resolver.EnsureTableExists(ctx, client, skKey)
		if err != nil {
			return nil, err
		}
	}

	return client, err
}
