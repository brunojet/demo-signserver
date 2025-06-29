package db_services

import (
	"context"
	"testing"

	db_mock "demo-signserver/pkg/repository/mock"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

func TestDynamoDBService_PutItem(t *testing.T) {
	mock := &db_mock.MockDynamoDBClient{
		PutItemFunc: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			assert.Equal(t, "TestTable", *params.TableName)
			assert.NotNil(t, params.Item["ID"])
			return &dynamodb.PutItemOutput{}, nil
		},
	}
	service := &DynamoDBService{Client: mock, Table: "TestTable"}
	item := map[string]types.AttributeValue{"ID": &types.AttributeValueMemberS{Value: "123"}}
	err := service.CreateItem(context.TODO(), item)
	assert.NoError(t, err)
}

func TestDynamoDBService_GetItem(t *testing.T) {
	mock := &db_mock.MockDynamoDBClient{
		GetItemFunc: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			assert.Equal(t, "TestTable", *params.TableName)
			assert.NotNil(t, params.Key["PK"])
			// Simula item retornado com PK (e SK se necessário)
			return &dynamodb.GetItemOutput{
				Item: map[string]types.AttributeValue{"PK": &types.AttributeValueMemberS{Value: "123"}},
			}, nil
		},
	}
	service := &DynamoDBService{Client: mock, Table: "TestTable", PKKey: "PK", SKKey: ""}
	item, err := service.GetItem(context.TODO(), "123")
	assert.NoError(t, err)
	assert.Equal(t, "123", item["ID"].(*types.AttributeValueMemberS).Value)
}
