package db_services

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// DynamoDBEndpointResolver é um endpoint resolver customizado e utilitário para ambiente local/dev.
type DynamoDBEndpointResolver struct {
	EndpointURL string
	TableName   string
}

// ResolveEndpoint implementa a interface de resolução de endpoint customizado.
func (r *DynamoDBEndpointResolver) ResolveEndpoint(service, region string, options ...interface{}) (aws.Endpoint, error) {
	fmt.Printf("[DynamoDBService] Usando endpoint customizado: %s\n", r.EndpointURL)
	return aws.Endpoint{
		URL:           r.EndpointURL,
		SigningRegion: "us-east-1",
	}, nil
}
