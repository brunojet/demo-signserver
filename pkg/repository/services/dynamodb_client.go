package db_services

import (
	"context"
	"demo-signserver/pkg/observability"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

var (
	client *dynamodb.Client
)

type resolverV2 struct{}

func (*resolverV2) ResolveEndpoint(ctx context.Context, params dynamodb.EndpointParameters) (
	smithyendpoints.Endpoint, error,
) {
	return dynamodb.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

func newDynamoDBClient(ctx context.Context) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg, func(optFns *dynamodb.Options) {
		if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
			optFns.BaseEndpoint = aws.String(endpoint)
			optFns.EndpointResolverV2 = &resolverV2{}
		}
		observability.LogInfo(
			"DynamoDB client information",
			map[string]interface{}{
				"region":   cfg.Region,
				"endpoint": optFns.BaseEndpoint,
			},
		)
	})
	return client, err
}

func GetDynamoDBCLient() *dynamodb.Client {
	if client == nil {
		var err error
		client, err = newDynamoDBClient(context.Background())
		if err != nil {
			fmt.Printf("Error initializing DynamoDB client: %v\n", err)
		}
	}
	return client
}
