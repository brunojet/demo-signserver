package db_services

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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

type ctxKey string

const CtxNoOverwriteKey ctxKey = "dynamodb_no_overwrite"

// PutItem insere um item na tabela DynamoDB.
func (s *DynamoDBService) putItemInternal(ctx context.Context, item map[string]types.AttributeValue) error {
	AddPKSKToItem(item, s.PKKey, s.SKKey)
	SetTimestamps(item, true)

	var condExpr *string
	var exprAttrNames map[string]string

	if v := ctx.Value(CtxNoOverwriteKey); v != nil {
		if noOverwrite, ok := v.(bool); ok && noOverwrite {
			cond, names := BuildNoOverwriteCondition(s.PKKey, s.SKKey)
			condExpr = &cond
			exprAttrNames = names
		}
	}

	_, err := s.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:                &s.Table,
		Item:                     item,
		ConditionExpression:      condExpr,
		ExpressionAttributeNames: exprAttrNames,
	})
	if err != nil {
		fmt.Printf("[DynamoDBService] Erro ao inserir item na tabela %s: %v\n", s.Table, err)
	}
	return err
}

// updateItemInternal executa o update no DynamoDB e loga erro se houver, usando receiver para acesso ao client e configs.
func (s *DynamoDBService) updateItemInternal(ctx context.Context, ID string, item map[string]types.AttributeValue) error {
	SetTimestamps(item, false)

	updateExpr, exprAttrNames, exprAttrValues, err := BuildUpdateExpressionFromAVMap(item)
	if err != nil {
		return err
	}

	_, err = s.Client.(*dynamodb.Client).UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 &s.Table,
		Key:                       MakeKeyByID(ID, s.PKKey, s.SKKey),
		UpdateExpression:          &updateExpr,
		ExpressionAttributeValues: exprAttrValues,
		ExpressionAttributeNames:  exprAttrNames,
	})
	if err != nil {
		fmt.Printf("[DynamoDBService] Erro ao atualizar item com ID %s na tabela %s: %v\n", ID, s.Table, err)
	}
	return err
}

// GetItem busca um item pela chave.
func (s *DynamoDBService) GetItem(ctx context.Context, ID string, out interface{}) error {
	resp, err := s.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &s.Table,
		Key:       MakeKeyByID(ID, s.PKKey, s.SKKey),
	})
	if err != nil {
		fmt.Printf("[DynamoDBService] Erro ao buscar item com ID %s na tabela %s: %v\n", ID, s.Table, err)
		return err
	}
	AddIDToItem(resp.Item, s.SKKey)
	return UnmarshalItem(resp.Item, out)
}

// CreateItem insere um item, sempre evitando sobrescrita (ConditionExpression).
func (s *DynamoDBService) CreateItem(ctx context.Context, obj interface{}) error {
	ctx = context.WithValue(ctx, CtxNoOverwriteKey, true)
	item, err := MarshalItem(obj)
	if err != nil {
		return err
	}
	return s.putItemInternal(ctx, item)
}

// UpdateItem atualiza apenas os campos não-chave do objeto informado (update parcial).
func (s *DynamoDBService) UpdateItem(ctx context.Context, ID string, obj interface{}) error {
	item, err := attributevalue.MarshalMap(obj)
	if err != nil {
		return err
	}
	return s.updateItemInternal(ctx, ID, item)
}
