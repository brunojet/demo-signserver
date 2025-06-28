package db_services

import (
	"reflect"
	"strings"
	"time"

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

// getReflectObj retorna o reflect.Value do objeto, já desreferenciado se for ponteiro (inclusive map)
func getReflectObj(obj interface{}) reflect.Value {
	v := reflect.ValueOf(obj)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

// getField retorna o campo (struct) ou valor (map) se existir, sem criar, e só retorna se for do tipo esperado
func getField(v reflect.Value, fieldName string, expectedKind reflect.Kind) reflect.Value {
	var val reflect.Value
	if v.Kind() == reflect.Struct {
		val = v.FieldByName(fieldName)
	} else if v.Kind() == reflect.Map {
		val = v.MapIndex(reflect.ValueOf(fieldName))
	}
	if val.IsValid() {
		valKind := val.Kind()

		// Se for interface, extrai o valor real
		if valKind == reflect.Interface && !val.IsNil() {
			val = val.Elem()
			valKind = val.Kind()
		}

		if valKind == expectedKind {
			return val
		}

		// Caso especial: map[string]interface{} com valor string esperado, mas valor não existe
		if valKind == reflect.Map && v.Type().Key().Kind() == reflect.String && v.Type().Elem().Kind() == reflect.Interface && expectedKind == reflect.String {
			return reflect.ValueOf("")
		}
	}
	return reflect.Value{}
}

// setField atribui valor string ao campo (struct ou map) se possível, criando no map se necessário
func setField(v reflect.Value, fieldName, value string) {
	if v.Kind() == reflect.Struct {
		f := v.FieldByName(fieldName)
		if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
			f.SetString(value)
		}
	} else if v.Kind() == reflect.Map {
		// Sempre adiciona ao map, mesmo se value for vazio
		v.SetMapIndex(reflect.ValueOf(fieldName), reflect.ValueOf(value))
	}
}

// getOrCreateField centraliza acesso e criação de campos em struct/map via reflection
func getOrCreateField(v reflect.Value, fieldName string, expectedKind reflect.Kind) reflect.Value {
	// Primeiro tenta obter o campo normalmente
	if f := getField(v, fieldName, expectedKind); f.IsValid() {
		return f
	}
	// Se não existir e for map, cria valor string vazio
	if v.Kind() == reflect.Map {
		key := reflect.ValueOf(fieldName)
		val := reflect.ValueOf("")
		v.SetMapIndex(key, val)
		return v.MapIndex(key)
	}
	return reflect.Value{}
}

// SetPKSKFromID preenche PK e SK a partir de ID (se existir). Se PK/SK não existirem, adiciona via reflection.
func SetPKSKFromID(obj interface{}) {
	v := getReflectObj(obj)
	idField := getField(v, "ID", reflect.String)
	if !idField.IsValid() || idField.String() == "" {
		return
	}
	id := idField.String()
	parts := strings.SplitN(id, ";", 2)
	pk, sk := parts[0], ""
	if len(parts) > 1 {
		sk = parts[1]
	}
	setField(v, "PK", pk)
	if sk != "" {
		setField(v, "SK", sk)
	}
}

// SetIDFromPKSK preenche ID a partir de PK e SK. Se ID não existir, adiciona via reflection (para map).
func SetIDFromPKSK(obj interface{}) {
	v := getReflectObj(obj)
	pkField := getField(v, "PK", reflect.String)
	skField := getField(v, "SK", reflect.String)
	var pk, sk string
	if pkField.IsValid() {
		pk = pkField.String()
	}

	if skField.IsValid() {
		sk = skField.String()
	}
	id := ""
	if pk != "" && sk != "" {
		id = pk + ";" + sk
	} else if pk != "" {
		id = pk
	}
	setField(v, "ID", id)
}

// SetTimestamps preenche CreatedAt e UpdatedAt se existirem no struct ou map
func SetTimestamps(obj interface{}, isCreate bool) {
	v := getReflectObj(obj)
	now := time.Now().Unix()
	if v.Kind() == reflect.Struct {
		updatedAt := getOrCreateField(v, "UpdatedAt", reflect.Int64)
		if isCreate {
			createdAt := getOrCreateField(v, "CreatedAt", reflect.Int64)
			if createdAt.IsValid() && createdAt.CanSet() {
				createdAt.SetInt(now)
			}
		}
		if updatedAt.IsValid() && updatedAt.CanSet() {
			updatedAt.SetInt(now)
		}
	} else if v.Kind() == reflect.Map {
		// Para map, setar diretamente
		if isCreate {
			v.SetMapIndex(reflect.ValueOf("CreatedAt"), reflect.ValueOf(now))
		}
		v.SetMapIndex(reflect.ValueOf("UpdatedAt"), reflect.ValueOf(now))
	}
}
