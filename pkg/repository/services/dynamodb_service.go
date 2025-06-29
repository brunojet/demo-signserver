package db_services

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
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
	PKKey  string
	SKKey  string
}

// NewDynamoDBServiceWithConfigLoader permite injetar função de carregamento de config (para testes).
func NewDynamoDBServiceWithConfigLoader(table string, pkKey string, skKey string, loadConfig func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error)) (*DynamoDBService, error) {
	cfg, err := loadConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	return &DynamoDBService{Client: client, Table: table, PKKey: pkKey, SKKey: skKey}, nil
}

// NewDynamoDBService padrão, usa config.LoadDefaultConfig
func NewDynamoDBService(table string, pkKey string, skKey string) (*DynamoDBService, error) {
	return NewDynamoDBServiceWithConfigLoader(table, pkKey, skKey, config.LoadDefaultConfig)
}

// PutItem insere um item na tabela DynamoDB.
func (s *DynamoDBService) PutItem(ctx context.Context, item map[string]types.AttributeValue) error {
	AddPKSKToItem(item, s.PKKey, s.SKKey)
	AddIDToItem(item, s.SKKey)
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
	AddIDToItem(resp.Item, s.SKKey)
	return resp.Item, nil
}

// getStringAttrValue extrai o valor string de um types.AttributeValueMemberS, ou retorna "" se não for string
func getStringAttrValue(attr types.AttributeValue) string {
	if v, ok := attr.(*types.AttributeValueMemberS); ok {
		return v.Value
	}
	return ""
}

// addPKSKToItem adiciona PK e SK ao item conforme as regras de negócio.
func (s *DynamoDBService) addPKSKToItem(item map[string]types.AttributeValue) {
	pkVal, pkOk := item[s.PKKey]

	if s.PKKey == "" || !pkOk {
		item["PK"] = &types.AttributeValueMemberS{Value: uuid.NewString()}
	} else {
		item["PK"] = &types.AttributeValueMemberS{Value: getStringAttrValue(pkVal)}
	}

	if s.SKKey != "" {
		skVal := item[s.SKKey]
		item["SK"] = &types.AttributeValueMemberS{Value: getStringAttrValue(skVal)}
	}
}

// addIDToItem preenche o campo ID no item a partir de PK e SK
func (s *DynamoDBService) addIDToItem(item map[string]types.AttributeValue) {
	pkStr := getStringAttrValue(item["PK"])
	if pkStr == "" {
		return
	}
	id := pkStr

	if s.SKKey != "" {
		skStr := getStringAttrValue(item["SK"])
		if skStr != "" {
			id += ";" + skStr
		}
	}

	item["ID"] = &types.AttributeValueMemberS{Value: id}

}
