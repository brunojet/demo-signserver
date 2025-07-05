package db_services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type TableServices struct {
	Client *dynamodb.Client
}

func NewDB() *TableServices {
	client := GetDynamoDBCLient()
	return &TableServices{Client: client}
}

func (db *TableServices) TableExists(ctx context.Context, table string) bool {
	_, err := db.Client.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: &table})
	if err == nil {
		return true
	}
	var rnfe *types.ResourceNotFoundException
	if errors.As(err, &rnfe) {
		return false
	}
	return false
}

// DeleteTable remove a tabela se ela existir (útil para testes).
func (db *TableServices) DeleteTable(ctx context.Context, table string) error {
	if !db.TableExists(ctx, table) {
		return nil
	}

	_, err := db.Client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: &table,
	})

	if err != nil {
		fmt.Printf("[DynamoDBService] Erro ao deletar tabela %s: %v\n", table, err)
		return fmt.Errorf("erro ao deletar tabela %s: %w", table, err)
	}

	fmt.Printf("[DynamoDBService] Tabela %s deletada com sucesso.\n", table)
	return nil
}

// CreateTable cria uma tabela DynamoDB com pk obrigatória e sk opcional (sempre usando nomes físicos pk/sk).
func (db *TableServices) CreateTable(ctx context.Context, table string, skKey string) error {
	if db.TableExists(ctx, table) {
		return nil
	}

	pkPhysical := PARTITION_KEY
	attrs := []types.AttributeDefinition{{AttributeName: &pkPhysical, AttributeType: types.ScalarAttributeTypeS}}
	keySchema := []types.KeySchemaElement{{AttributeName: &pkPhysical, KeyType: types.KeyTypeHash}}

	if skKey != "" {
		skPhysical := SORT_KEY
		attrs = append(attrs, types.AttributeDefinition{AttributeName: &skPhysical, AttributeType: types.ScalarAttributeTypeS})
		keySchema = append(keySchema, types.KeySchemaElement{AttributeName: &skPhysical, KeyType: types.KeyTypeRange})
	}

	_, err := db.Client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName:            &table,
		AttributeDefinitions: attrs,
		KeySchema:            keySchema,
		BillingMode:          types.BillingModePayPerRequest,
	})
	if err != nil {
		fmt.Printf("[DynamoDBService] Erro ao criar tabela %s: %v\n", table, err)
		return fmt.Errorf("erro ao criar tabela %s: %w", table, err)
	}
	fmt.Printf("[DynamoDBService] Tabela %s criada com sucesso (pk/sk).\n", table)
	return nil
}

func (db *TableServices) CreateRandomTable(ctx context.Context, skKey string) string {
	tableName := uuid.New().String()
	err := db.CreateTable(ctx, tableName, skKey)
	if err != nil {
		log.Fatalf("erro ao criar tabela aleatória %s: %v", tableName, err)
		return ""
	}
	fmt.Printf("[DynamoDBService] Tabela aleatória %s criada com sucesso.\n", tableName)
	return tableName
}
