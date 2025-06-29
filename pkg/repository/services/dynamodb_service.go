package db_services

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDBService encapsula operações com o DynamoDB.
type DynamoDBService struct {
	Client *dynamodb.Client
	Table  string
	PKKey  string
	SKKey  string
}

// NewDynamoDBService padrão, usa config.LoadDefaultConfig
func NewDynamoDBService(table string, pkKey string, skKey string) (*DynamoDBService, error) {
	project := os.Getenv("PROJECT_NAME")
	env := os.Getenv("ENVIRONMENT")
	table_name := fmt.Sprintf("%s-%s-%s", project, env, table)
	client, resolver, err := NewDynamoDBClient(context.TODO(), table_name)

	if err != nil {
		return nil, err
	}

	// Se for endpoint customizado, garante que a tabela existe
	if resolver.EndpointURL != "" {
		err = resolver.EnsureTableExists(context.TODO(), client, skKey)
		if err != nil {
			return nil, err
		}
	}

	return &DynamoDBService{Client: client, Table: table, PKKey: pkKey, SKKey: skKey}, nil
}

// PutItem insere um item na tabela DynamoDB.
func (s *DynamoDBService) putItemInternal(ctx context.Context, item map[string]types.AttributeValue) error {
	AddPKSKToItem(item, s.PKKey, s.SKKey)
	SetTimestamps(item, true)

	condExpr, exprAttrNames := BuildNoOverwriteCondition(s.PKKey, s.SKKey)

	_, err := s.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:                &s.Table,
		Item:                     item,
		ConditionExpression:      &condExpr,
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

	condExpr := BuildUpdateCondition(s.PKKey, s.SKKey)

	_, err = s.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 &s.Table,
		Key:                       MakeKeyByID(ID, s.PKKey, s.SKKey),
		UpdateExpression:          &updateExpr,
		ExpressionAttributeValues: exprAttrValues,
		ExpressionAttributeNames:  exprAttrNames,
		ConditionExpression:       &condExpr,
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
	if err != nil || len(resp.Item) == 0 {
		if err == nil {
			err = fmt.Errorf("[DynamoDBService] Erro ao buscar item com ID %s na tabela %s", ID, s.Table)
		}
		return err
	}
	AddIDToItem(resp.Item, s.SKKey)
	return UnmarshalItem(resp.Item, out)
}

// CreateItem insere um item, sempre evitando sobrescrita (ConditionExpression).
func (s *DynamoDBService) CreateItem(ctx context.Context, obj interface{}) error {
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
