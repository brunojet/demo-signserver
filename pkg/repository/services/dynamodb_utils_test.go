package db_services

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Mock para BaseDomainInterface

type mockBaseDomain struct {
	id string
}

func (m *mockBaseDomain) GetID() string   { return m.id }
func (m *mockBaseDomain) SetID(id string) { m.id = id }
func (m *mockBaseDomain) SetCreateTs()    {}
func (m *mockBaseDomain) SetUpdateTs()    {}

func TestAddPKSKToItem_OnlyPK(t *testing.T) {
	item := map[string]types.AttributeValue{
		"User": &types.AttributeValueMemberS{Value: "user1"},
	}
	obj := &mockBaseDomain{id: "user1"}
	AddPKSKToItem(obj, item, "User", "")
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
	obj := &mockBaseDomain{id: "user1"}
	AddPKSKToItem(obj, item, "User", "Type")
	if v, ok := item[PARTITION_KEY].(*types.AttributeValueMemberS); !ok || v.Value != "user1" {
		t.Errorf("PK not set correctly, got %v", item[PARTITION_KEY])
	}
	if v, ok := item[SORT_KEY].(*types.AttributeValueMemberS); !ok || v.Value != "admin" {
		t.Errorf("SK not set correctly, got %v", item[SORT_KEY])
	}
}

func TestAddPKSKToItem_NoPKKey(t *testing.T) {
	item := map[string]types.AttributeValue{}
	obj := &mockBaseDomain{id: ""}
	AddPKSKToItem(obj, item, "", "Type")
	if v, ok := item[PARTITION_KEY].(*types.AttributeValueMemberS); !ok || v.Value == "" {
		t.Errorf("PK should be a generated UUID, got %v", item[PARTITION_KEY])
	}
}

func TestAddIDToItem_OnlyPK(t *testing.T) {
	item := map[string]types.AttributeValue{
		PARTITION_KEY: &types.AttributeValueMemberS{Value: "user1"},
	}
	obj := &mockBaseDomain{id: "user1"}
	AddIDToItem(obj, item, ID_KEY, "")
	if v, ok := item[ID_KEY].(*types.AttributeValueMemberS); !ok || v.Value != "user1" {
		t.Errorf("ID not set correctly, got %v", item[ID_KEY])
	}
}

func TestAddIDToItem_PKSK(t *testing.T) {
	item := map[string]types.AttributeValue{
		PARTITION_KEY: &types.AttributeValueMemberS{Value: "user1"},
		SORT_KEY:      &types.AttributeValueMemberS{Value: "admin"},
	}
	obj := &mockBaseDomain{id: "user1"}
	AddIDToItem(obj, item, PARTITION_KEY, SORT_KEY)
	if v, ok := item[ID_KEY].(*types.AttributeValueMemberS); !ok || v.Value != "user1-admin" {
		t.Errorf("ID not set correctly, got %v", item[ID_KEY])
	}
}

func TestAddIDToItem_NoPK(t *testing.T) {
	item := map[string]types.AttributeValue{
		SORT_KEY: &types.AttributeValueMemberS{Value: "admin"},
	}
	obj := &mockBaseDomain{id: ""}
	AddIDToItem(obj, item, PARTITION_KEY, SORT_KEY)
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
	key := "pkval-skval"
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

func TestBuildUpdateExpression(t *testing.T) {
	fields := map[string]interface{}{
		"foo": "bar",
		"num": 42,
	}
	expr, names, values, err := BuildUpdateExpression(fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if expr != "SET #foo = :foo, #num = :num" && expr != "SET #num = :num, #foo = :foo" {
		t.Errorf("unexpected expr: %s", expr)
	}
	if names["#foo"] != "foo" || names["#num"] != "num" {
		t.Errorf("unexpected names: %#v", names)
	}
	if _, ok := values[":foo"]; !ok {
		t.Errorf(":foo missing in values")
	}
	if _, ok := values[":num"]; !ok {
		t.Errorf(":num missing in values")
	}
}

func TestIsNonUpdatable(t *testing.T) {
	for _, k := range NonUpdatableKeys {
		if !isNonUpdatable(k) {
			t.Errorf("expected %s to be non-updatable", k)
		}
	}
	if isNonUpdatable("other") {
		t.Errorf("expected 'other' to be updatable")
	}
}

func TestBuildUpdateExpressionFromAVMap(t *testing.T) {
	b := map[string]types.AttributeValue{
		"foo":         &types.AttributeValueMemberS{Value: "bar"},
		"num":         &types.AttributeValueMemberN{Value: "42"},
		PARTITION_KEY: &types.AttributeValueMemberS{Value: "pkval"},
		ID_KEY:        &types.AttributeValueMemberS{Value: "idval"},
	}
	expr, names, values, err := BuildUpdateExpressionFromAVMap(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if expr != "SET #foo = :foo, #num = :num" && expr != "SET #num = :num, #foo = :foo" {
		t.Errorf("unexpected expr: %s", expr)
	}
	if names["#foo"] != "foo" || names["#num"] != "num" {
		t.Errorf("unexpected names: %#v", names)
	}
	if _, ok := values[":foo"]; !ok {
		t.Errorf(":foo missing in values")
	}
	if _, ok := values[":num"]; !ok {
		t.Errorf(":num missing in values")
	}
	// Non-updatable keys should not be present
	for _, k := range []string{PARTITION_KEY, ID_KEY} {
		for n := range names {
			if names[n] == k {
				t.Errorf("non-updatable key %s present in names", k)
			}
		}
	}
}

func TestBuildNoOverwriteCondition(t *testing.T) {
	cond, names := BuildNoOverwriteCondition(ID_KEY, "type")
	if cond != "attribute_not_exists(#pk) AND attribute_not_exists(#sk)" {
		t.Errorf("unexpected cond: %s", cond)
	}
	if names["#pk"] != PARTITION_KEY || names["#sk"] != "type" {
		t.Errorf("unexpected names: %#v", names)
	}

	cond, names = BuildNoOverwriteCondition("", "")
	if cond != "attribute_not_exists(#pk)" {
		t.Errorf("unexpected cond for only PK: %s", cond)
	}
	// Aceita PARTITION_KEY ou "" como valor válido para #pk
	if v, ok := names["#pk"]; !ok || (v != PARTITION_KEY && v != "") {
		t.Errorf("unexpected pk name: %s", v)
	}
	if _, ok := names["#sk"]; ok {
		t.Errorf("should not have #sk in names")
	}
}
