package db_services

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestAddPKSKToItem_OnlyPK(t *testing.T) {
	item := map[string]types.AttributeValue{
		"User": &types.AttributeValueMemberS{Value: "user1"},
	}
	AddPKSKToItem(item, "User", "")
	if v, ok := item[PARTITION_KEY].(*types.AttributeValueMemberS); !ok || v.Value != "user1" {
		t.Errorf("PK not set correctly, got %v", item[PARTITION_KEY])
	}
	if _, ok := item[SORT_KEY]; ok {
		t.Errorf("SK should not be set when SKKey is empty")
	}
}

func TestAddPKSKToItem_PKSK(t *testing.T) {
	item := map[string]types.AttributeValue{
		"User": &types.AttributeValueMemberS{Value: "user1"},
		"Type": &types.AttributeValueMemberS{Value: "admin"},
	}
	AddPKSKToItem(item, "User", "Type")
	if v, ok := item[PARTITION_KEY].(*types.AttributeValueMemberS); !ok || v.Value != "user1" {
		t.Errorf("PK not set correctly, got %v", item[PARTITION_KEY])
	}
	if v, ok := item[SORT_KEY].(*types.AttributeValueMemberS); !ok || v.Value != "admin" {
		t.Errorf("SK not set correctly, got %v", item[SORT_KEY])
	}
}

func TestAddPKSKToItem_NoPKKey(t *testing.T) {
	item := map[string]types.AttributeValue{}
	AddPKSKToItem(item, "", "Type")
	if v, ok := item[PARTITION_KEY].(*types.AttributeValueMemberS); !ok || v.Value == "" {
		t.Errorf("PK should be a generated UUID, got %v", item[PARTITION_KEY])
	}
}

func TestAddIDToItem_OnlyPK(t *testing.T) {
	item := map[string]types.AttributeValue{
		PARTITION_KEY: &types.AttributeValueMemberS{Value: "user1"},
	}
	AddIDToItem(item, "")
	if v, ok := item[ID_KEY].(*types.AttributeValueMemberS); !ok || v.Value != "user1" {
		t.Errorf("ID not set correctly, got %v", item[ID_KEY])
	}
}

func TestAddIDToItem_PKSK(t *testing.T) {
	item := map[string]types.AttributeValue{
		PARTITION_KEY: &types.AttributeValueMemberS{Value: "user1"},
		SORT_KEY:      &types.AttributeValueMemberS{Value: "admin"},
	}
	AddIDToItem(item, SORT_KEY)
	if v, ok := item[ID_KEY].(*types.AttributeValueMemberS); !ok || v.Value != "user1#admin" {
		t.Errorf("ID not set correctly, got %v", item[ID_KEY])
	}
}

func TestAddIDToItem_NoPK(t *testing.T) {
	item := map[string]types.AttributeValue{
		SORT_KEY: &types.AttributeValueMemberS{Value: "admin"},
	}
	AddIDToItem(item, SORT_KEY)
	if _, ok := item[ID_KEY]; ok {
		t.Errorf("ID should not be set when PK is missing")
	}
}

func TestMakeKeyByID_OnlyPK(t *testing.T) {
	key := "pkvalue"
	result := MakeKeyByID(key, PARTITION_KEY, "")
	if v, ok := result[PARTITION_KEY].(*types.AttributeValueMemberS); !ok || v.Value != "pkvalue" {
		t.Errorf("Expected PK=pkvalue, got %v", result[PARTITION_KEY])
	}
	if _, ok := result[SORT_KEY]; ok {
		t.Errorf("SK should not be set when skKey is empty")
	}
}

func TestMakeKeyByID_PKSK(t *testing.T) {
	key := "pkval#skval"
	result := MakeKeyByID(key, PARTITION_KEY, SORT_KEY)
	if v, ok := result[PARTITION_KEY].(*types.AttributeValueMemberS); !ok || v.Value != "pkval" {
		t.Errorf("Expected PK=pkval, got %v", result[PARTITION_KEY])
	}
	if v, ok := result[SORT_KEY].(*types.AttributeValueMemberS); !ok || v.Value != "skval" {
		t.Errorf("Expected SK=skval, got %v", result[SORT_KEY])
	}
}

func TestMakeKeyByID_PKSK_OnlyPKProvided(t *testing.T) {
	key := "pkval"
	result := MakeKeyByID(key, PARTITION_KEY, SORT_KEY)
	if v, ok := result[PARTITION_KEY].(*types.AttributeValueMemberS); !ok || v.Value != "pkval" {
		t.Errorf("Expected PK=pkval, got %v", result[PARTITION_KEY])
	}
	if _, ok := result[SORT_KEY]; ok {
		t.Errorf("SK should not be set when only PK is provided")
	}
}

func TestMakeKeyByID_EmptyKey(t *testing.T) {
	result := MakeKeyByID("", PARTITION_KEY, SORT_KEY)
	if v, ok := result[PARTITION_KEY].(*types.AttributeValueMemberS); !ok || v.Value != "" {
		t.Errorf("Expected PK='', got %v", result[PARTITION_KEY])
	}
	if _, ok := result[SORT_KEY]; ok {
		t.Errorf("SK should not be set when key is empty")
	}
}
