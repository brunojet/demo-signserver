package db_services

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

type resolverV2 struct{}

func (*resolverV2) ResolveEndpoint(ctx context.Context, params dynamodb.EndpointParameters) (
	smithyendpoints.Endpoint, error,
) {
	// s3.Options.BaseEndpoint is accessible here:
	fmt.Printf("The endpoint provided in config is %s\n", *params.Endpoint)

	// fallback to default
	return dynamodb.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

func NewDynamoDBClient(ctx context.Context, table string) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg, func(optFns *dynamodb.Options) {
		if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
			optFns.BaseEndpoint = aws.String(endpoint)
			optFns.EndpointResolverV2 = &resolverV2{}
		}
	})

	return client, err
}
