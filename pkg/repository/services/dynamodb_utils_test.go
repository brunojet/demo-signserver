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
	if v, ok := item["PK"].(*types.AttributeValueMemberS); !ok || v.Value != "user1" {
		t.Errorf("PK not set correctly, got %v", item["PK"])
	}
	if _, ok := item["SK"]; ok {
		t.Errorf("SK should not be set when SKKey is empty")
	}
}

func TestAddPKSKToItem_PKSK(t *testing.T) {
	item := map[string]types.AttributeValue{
		"User": &types.AttributeValueMemberS{Value: "user1"},
		"Type": &types.AttributeValueMemberS{Value: "admin"},
	}
	AddPKSKToItem(item, "User", "Type")
	if v, ok := item["PK"].(*types.AttributeValueMemberS); !ok || v.Value != "user1" {
		t.Errorf("PK not set correctly, got %v", item["PK"])
	}
	if v, ok := item["SK"].(*types.AttributeValueMemberS); !ok || v.Value != "admin" {
		t.Errorf("SK not set correctly, got %v", item["SK"])
	}
}

func TestAddPKSKToItem_NoPKKey(t *testing.T) {
	item := map[string]types.AttributeValue{}
	AddPKSKToItem(item, "", "Type")
	if v, ok := item["PK"].(*types.AttributeValueMemberS); !ok || v.Value == "" {
		t.Errorf("PK should be a generated UUID, got %v", item["PK"])
	}
}

func TestAddIDToItem_OnlyPK(t *testing.T) {
	item := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: "user1"},
	}
	AddIDToItem(item, "")
	if v, ok := item["ID"].(*types.AttributeValueMemberS); !ok || v.Value != "user1" {
		t.Errorf("ID not set correctly, got %v", item["ID"])
	}
}

func TestAddIDToItem_PKSK(t *testing.T) {
	item := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: "user1"},
		"SK": &types.AttributeValueMemberS{Value: "admin"},
	}
	AddIDToItem(item, "SK")
	if v, ok := item["ID"].(*types.AttributeValueMemberS); !ok || v.Value != "user1#admin" {
		t.Errorf("ID not set correctly, got %v", item["ID"])
	}
}

func TestAddIDToItem_NoPK(t *testing.T) {
	item := map[string]types.AttributeValue{
		"SK": &types.AttributeValueMemberS{Value: "admin"},
	}
	AddIDToItem(item, "SK")
	if _, ok := item["ID"]; ok {
		t.Errorf("ID should not be set when PK is missing")
	}
}

func TestMakeKeyByID_OnlyPK(t *testing.T) {
	key := "pkvalue"
	result := MakeKeyByID(key, "PK", "")
	if v, ok := result["PK"].(*types.AttributeValueMemberS); !ok || v.Value != "pkvalue" {
		t.Errorf("Expected PK=pkvalue, got %v", result["PK"])
	}
	if _, ok := result["SK"]; ok {
		t.Errorf("SK should not be set when skKey is empty")
	}
}

func TestMakeKeyByID_PKSK(t *testing.T) {
	key := "pkval#skval"
	result := MakeKeyByID(key, "PK", "SK")
	if v, ok := result["PK"].(*types.AttributeValueMemberS); !ok || v.Value != "pkval" {
		t.Errorf("Expected PK=pkval, got %v", result["PK"])
	}
	if v, ok := result["SK"].(*types.AttributeValueMemberS); !ok || v.Value != "skval" {
		t.Errorf("Expected SK=skval, got %v", result["SK"])
	}
}

func TestMakeKeyByID_PKSK_OnlyPKProvided(t *testing.T) {
	key := "pkval"
	result := MakeKeyByID(key, "PK", "SK")
	if v, ok := result["PK"].(*types.AttributeValueMemberS); !ok || v.Value != "pkval" {
		t.Errorf("Expected PK=pkval, got %v", result["PK"])
	}
	if _, ok := result["SK"]; ok {
		t.Errorf("SK should not be set when only PK is provided")
	}
}

func TestMakeKeyByID_EmptyKey(t *testing.T) {
	result := MakeKeyByID("", "PK", "SK")
	if v, ok := result["PK"].(*types.AttributeValueMemberS); !ok || v.Value != "" {
		t.Errorf("Expected PK='', got %v", result["PK"])
	}
	if _, ok := result["SK"]; ok {
		t.Errorf("SK should not be set when key is empty")
	}
}
