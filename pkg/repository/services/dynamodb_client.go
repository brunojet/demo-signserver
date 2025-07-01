package db_services

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewDynamoDBClient(ctx context.Context, table string) (*dynamodb.Client, error) {
	loadConfig := config.LoadDefaultConfig
	cfg, err := loadConfig(ctx, func(o *config.LoadOptions) error {
		if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
			resolver := &DynamoDBEndpointResolver{
				EndpointURL: endpoint,
				TableName:   table,
			}
			o.EndpointResolverWithOptions = resolver
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg)

	return client, err
}
