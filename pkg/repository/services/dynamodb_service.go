package db_services

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDBAPI define interface para mocks do DynamoDB Client.
type DynamoDBAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

// DynamoDBService encapsula operações com o DynamoDB.
type DynamoDBService struct {
	Client DynamoDBAPI
	Table  string
}

// NewDynamoDBServiceWithConfigLoader permite injetar função de carregamento de config (para testes).
func NewDynamoDBServiceWithConfigLoader(table string, loadConfig func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error)) (*DynamoDBService, error) {
	cfg, err := loadConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	return &DynamoDBService{Client: client, Table: table}, nil
}

// NewDynamoDBService padrão, usa config.LoadDefaultConfig
func NewDynamoDBService(table string) (*DynamoDBService, error) {
	return NewDynamoDBServiceWithConfigLoader(table, config.LoadDefaultConfig)
}

// PutItem insere um item na tabela DynamoDB.
func (s *DynamoDBService) PutItem(ctx context.Context, item map[string]types.AttributeValue) error {
	_, err := s.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &s.Table,
		Item:      item,
	})
	return err
}

// GetItem busca um item pela chave.
func (s *DynamoDBService) GetItem(ctx context.Context, key map[string]types.AttributeValue) (map[string]types.AttributeValue, error) {
	resp, err := s.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &s.Table,
		Key:       key,
	})
	if err != nil {
		return nil, err
	}
	return resp.Item, nil
}
