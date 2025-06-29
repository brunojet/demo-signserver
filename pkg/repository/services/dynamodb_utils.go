package db_services

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

const PARTITION_KEY = "pk"
const SORT_KEY = "sk"
const ID_KEY = "id"
const CREATED_AT_KEY = "created_at"
const UPDATED_AT_KEY = "updated_at"

// MarshalItem converte struct Go para map[string]types.AttributeValue (DynamoDB)
func MarshalItem(v interface{}) (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(v)
}

// UnmarshalItem converte map[string]types.AttributeValue para struct Go (DynamoDB)
func UnmarshalItem(m map[string]types.AttributeValue, out interface{}) error {
	return attributevalue.UnmarshalMap(m, out)
}

// setTimestamps preenche CreatedAt e UpdatedAt em item map[string]types.AttributeValue
func SetTimestamps(item map[string]types.AttributeValue, isCreate bool) {
	now := time.Now().Unix()
	if isCreate {
		item[CREATED_AT_KEY] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", now)}
	}
	item[UPDATED_AT_KEY] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", now)}
}

// getStringAttrValue extrai o valor string de um types.AttributeValueMemberS, ou retorna "" se não for string
func GetStringAttrValue(attr types.AttributeValue) string {
	if v, ok := attr.(*types.AttributeValueMemberS); ok {
		return v.Value
	}
	return ""
}

// AddPKSKToItem adiciona PK e SK ao item conforme as regras de negócio.
func AddPKSKToItem(item map[string]types.AttributeValue, pkKey, skKey string) {
	pkVal, pkOk := item[pkKey]

	if pkKey == "" || !pkOk {
		item[PARTITION_KEY] = &types.AttributeValueMemberS{Value: uuid.NewString()}
	} else {
		item[PARTITION_KEY] = &types.AttributeValueMemberS{Value: GetStringAttrValue(pkVal)}
	}

	if skKey != "" {
		skVal := item[skKey]
		item[SORT_KEY] = &types.AttributeValueMemberS{Value: GetStringAttrValue(skVal)}
	}
}

// AddIDToItem preenche o campo ID no item a partir de PK e SK
func AddIDToItem(item map[string]types.AttributeValue, skKey string) {
	pkStr := GetStringAttrValue(item[PARTITION_KEY])
	if pkStr == "" {
		return
	}
	id := pkStr

	if skKey != "" {
		skStr := GetStringAttrValue(item[SORT_KEY])
		if skStr != "" {
			id += "#" + skStr
		}
	}

	item[ID_KEY] = &types.AttributeValueMemberS{Value: id}
}

// MakeKeyByID cria a chave para busca no DynamoDB: PK ou PK/SK
func MakeKeyByID(key, pkKey, skKey string) map[string]types.AttributeValue {
	result := make(map[string]types.AttributeValue)
	if skKey != "" {
		parts := strings.SplitN(key, "#", 2)
		result[PARTITION_KEY] = &types.AttributeValueMemberS{Value: parts[0]}
		if len(parts) > 1 {
			result[SORT_KEY] = &types.AttributeValueMemberS{Value: parts[1]}
		}
	} else {
		result[PARTITION_KEY] = &types.AttributeValueMemberS{Value: key}
	}

	return result
}
