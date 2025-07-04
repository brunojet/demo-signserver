package db_services

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DBServices struct {
	Client *dynamodb.Client
}

func NewDB(table string) *DBServices {
	client := GetDynamoDBCLient()
	return &DBServices{Client: client}
}

func buildTableName(table string) string {
	project := os.Getenv("PROJECT_NAME")
	env := os.Getenv("ENVIRONMENT")
	return fmt.Sprintf("%s-%s-%s", project, env, table)
}

// TableExists retorna true se a tabela existe, false se não existe, ou erro se outro erro.
func (db *DBServices) tableExists(ctx context.Context, table string) bool {
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

func (db *DBServices) TableExists(ctx context.Context, table string) bool {
	table = buildTableName(table)
	return db.tableExists(ctx, table)
}

// DeleteTable remove a tabela se ela existir (útil para testes).
func (db *DBServices) DeleteTable(ctx context.Context, table string) error {
	table = buildTableName(table)

	if !db.tableExists(ctx, table) {
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
func (db *DBServices) CreateTable(ctx context.Context, table string, skKey string) error {
	table = buildTableName(table)

	if db.tableExists(ctx, table) {
		return nil
	}

	pkPhysical := "pk"
	attrs := []types.AttributeDefinition{{AttributeName: &pkPhysical, AttributeType: types.ScalarAttributeTypeS}}
	keySchema := []types.KeySchemaElement{{AttributeName: &pkPhysical, KeyType: types.KeyTypeHash}}

	if skKey != "" {
		skPhysical := "sk"
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
