package db_services

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// MarshalItem converte struct Go para map[string]types.AttributeValue (DynamoDB)
func MarshalItem(v interface{}) (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(v)
}

// UnmarshalItem converte map[string]types.AttributeValue para struct Go (DynamoDB)
func UnmarshalItem(m map[string]types.AttributeValue, out interface{}) error {
	return attributevalue.UnmarshalMap(m, out)
}

// BuildKeyString constrói uma chave simples para consultas por PK string
func BuildKeyString(pkName, pkValue string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		pkName: &types.AttributeValueMemberS{Value: pkValue},
	}
}

// MarshalKey converte um map[string]interface{} para map[string]types.AttributeValue
func MarshalKey(key map[string]interface{}) map[string]types.AttributeValue {
	result := make(map[string]types.AttributeValue)
	for k, v := range key {
		if s, ok := v.(string); ok {
			result[k] = &types.AttributeValueMemberS{Value: s}
		}
		// Adicione outros tipos conforme necessário
	}
	return result
}

// BuildKeyFromID recebe id no formato "pk" ou "pk;sk" e retorna map[string]types.AttributeValue
func BuildKeyFromID(id string) map[string]types.AttributeValue {
	pk, sk := id, ""
	if strings.Contains(id, ";") {
		parts := strings.SplitN(id, ";", 2)
		pk = parts[0]
		sk = parts[1]
	}
	key := map[string]types.AttributeValue{"pk": &types.AttributeValueMemberS{Value: pk}}
	if sk != "" {
		key["sk"] = &types.AttributeValueMemberS{Value: sk}
	}
	return key
}
