package db_services

import (
	"context"
	"demo-signserver/pkg/repository/domain"
	"fmt"
	"log"

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
func NewDynamoDBService(table string, pkKey string, skKey string) *DynamoDBService {
	client := GetDynamoDBCLient()
	if pkKey == "" || pkKey == PARTITION_KEY {
		log.Fatalf("[DynamoDBService] pkKey não pode ser vazio ou igual a '%s'", PARTITION_KEY)
	} else if skKey == SORT_KEY {
		log.Fatalf("[DynamoDBService] skKey não pode ser igual a '%s'", SORT_KEY)
	}
	return &DynamoDBService{Client: client, Table: table, PKKey: pkKey, SKKey: skKey}
}

// PutItem insere um item na tabela DynamoDB.
func (s *DynamoDBService) putItemInternal(ctx context.Context, item map[string]types.AttributeValue) error {
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

// HasID define interface para structs que possuem ID
// Deve ser implementada por todos que embutem BaseDomain
// Exemplo: func (s *SignRequest) SetID(id string) { s.ID = id }

// GetItem busca um item pela chave.
func (s *DynamoDBService) GetItem(ctx context.Context, ID string, out domain.BaseDomainInterface) error {
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

	err = UnmarshalItem(resp.Item, out)
	if err != nil {
		return err
	}

	out.SetID(ID)

	return nil
}

// CreateItem insere um item, sempre evitando sobrescrita (ConditionExpression).
func (s *DynamoDBService) CreateItem(ctx context.Context, obj domain.BaseDomainInterface) error {
	obj.SetCreateTs()
	item, err := MarshalItem(obj)
	if err != nil {
		return err
	}

	AddPKSKToItem(obj, item, s.PKKey, s.SKKey)
	err = s.putItemInternal(ctx, item)
	AddIDToItem(obj, item, s.PKKey, s.SKKey)

	if err != nil {
		return err
	}

	obj.SetID(GetStringAttrValue(item[ID_KEY]))

	return err
}

// UpdateItem atualiza apenas os campos não-chave do objeto informado (update parcial).
func (s *DynamoDBService) UpdateItem(ctx context.Context, ID string, obj domain.BaseDomainInterface) error {
	obj.SetUpdateTs()
	item, err := attributevalue.MarshalMap(obj)

	if err != nil {
		return err
	}

	return s.updateItemInternal(ctx, ID, item)
}
