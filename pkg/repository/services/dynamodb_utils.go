package db_services

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

const (
	PARTITION_KEY  = "pk"
	SORT_KEY       = "sk"
	ID_KEY         = "id"
	CREATED_AT_KEY = "created_at"
	UPDATED_AT_KEY = "updated_at"
)

var NonUpdatableKeys = []string{PARTITION_KEY, SORT_KEY, ID_KEY, CREATED_AT_KEY}

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
func MakeKeyByID(ID, pkKey, skKey string) map[string]types.AttributeValue {
	result := make(map[string]types.AttributeValue)
	if skKey != "" {
		parts := strings.SplitN(ID, "#", 2)
		result[PARTITION_KEY] = &types.AttributeValueMemberS{Value: parts[0]}
		if len(parts) > 1 {
			result[SORT_KEY] = &types.AttributeValueMemberS{Value: parts[1]}
		}
	} else {
		result[PARTITION_KEY] = &types.AttributeValueMemberS{Value: ID}
	}

	return result
}

// BuildUpdateExpression monta a UpdateExpression, ExpressionAttributeNames e ExpressionAttributeValues para update parcial no DynamoDB.
// Ignora campos com valor nil.
func BuildUpdateExpression(fields map[string]interface{}) (string, map[string]string, map[string]types.AttributeValue, error) {
	updateExpr := "SET "
	exprAttrValues := make(map[string]types.AttributeValue)
	exprAttrNames := make(map[string]string)
	first := true
	for k, v := range fields {
		if v == nil {
			continue // ignora campos nil
		}
		if !first {
			updateExpr += ", "
		}
		first = false
		phName := "#" + k
		phValue := ":" + k
		updateExpr += phName + " = " + phValue
		exprAttrNames[phName] = k
		av, err := attributevalue.Marshal(v)
		if err != nil {
			return "", nil, nil, err
		}
		exprAttrValues[phValue] = av
	}
	if first { // nenhum campo válido
		return "", nil, nil, nil
	}
	return updateExpr, exprAttrNames, exprAttrValues, nil
}

func isNonUpdatable(key string) bool {
	for _, nonUpdatableKey := range NonUpdatableKeys {
		if key == nonUpdatableKey {
			return true
		}
	}
	return false
}

// BuildUpdateExpressionFromAVMap recebe um map[string]types.AttributeValue (MarshalMap) e monta a UpdateExpression ignorando PK, SK, created_at, updated_at e id.
// Ignora também campos com valor nil.
func BuildUpdateExpressionFromAVMap(b map[string]types.AttributeValue) (string, map[string]string, map[string]types.AttributeValue, error) {
	fields := make(map[string]interface{})
	for k, v := range b {
		if isNonUpdatable(k) || v == nil {
			continue
		}
		var val interface{}
		_ = attributevalue.Unmarshal(v, &val)
		if val == nil {
			continue
		}
		fields[k] = val
	}
	return BuildUpdateExpression(fields)
}

// BuildNoOverwriteCondition monta a ConditionExpression e ExpressionAttributeNames para evitar sobrescrita de item no DynamoDB.
func BuildNoOverwriteCondition(pkKey, skKey string) (condExpr string, exprAttrNames map[string]string) {
	pkName := pkKey
	if pkName == "" {
		pkName = PARTITION_KEY
	}
	cond := "attribute_not_exists(#pk)"
	exprAttrNames = map[string]string{"#pk": pkName}
	if skKey != "" {
		skName := skKey
		cond += " AND attribute_not_exists(#sk)"
		exprAttrNames["#sk"] = skName
	}
	return cond, exprAttrNames
}

// BuildUpdateCondition retorna uma ConditionExpression para garantir que o item existe antes do update.
func BuildUpdateCondition(pkKey, skKey string) string {
	cond := fmt.Sprintf("attribute_exists(%s)", pkKey)
	if skKey != "" {
		cond += fmt.Sprintf(" AND attribute_exists(%s)", skKey)
	}
	return cond
}
