package db_services

import (
	"reflect"
	"testing"
	"time"
)

type testStructWithId struct {
	ID string
}

type testStructWithPk struct {
	PK string
}

type testStructWithPkSk struct {
	PK string
	SK string
}

type testStruct struct {
	ID        string
	PK        string
	SK        string
	CreatedAt int64
	UpdatedAt int64
}

func TestSetPKSKFromID_Struct(t *testing.T) {
	ts := &testStruct{ID: "abc;def"}
	SetPKSKFromID(ts)
	if ts.PK != "abc" || ts.SK != "def" {
		t.Errorf("PK/SK not set correctly: got PK=%s, SK=%s", ts.PK, ts.SK)
	}
}

func TestSetPKSKFromID_Map(t *testing.T) {
	m := map[string]interface{}{"ID": "x;y"}
	SetPKSKFromID(&m)
	if m["PK"] != "x" || m["SK"] != "y" {
		t.Errorf("PK/SK not set correctly in map: got PK=%v, SK=%v", m["PK"], m["SK"])
	}
}

func TestSetIDFromPKSK_Struct(t *testing.T) {
	ts := &testStruct{PK: "a", SK: "b"}
	SetIDFromPKSK(ts)
	if ts.ID != "a;b" {
		t.Errorf("ID not set correctly: got %s", ts.ID)
	}
}

func TestSetIDFromPKSK_Map(t *testing.T) {
	m := map[string]interface{}{"PK": "foo", "SK": "bar"}
	SetIDFromPKSK(&m)
	if m["ID"] != "foo;bar" {
		t.Errorf("ID not set correctly in map: got %v", m["ID"])
	}
}

func TestSetTimestamps_Struct(t *testing.T) {
	ts := &testStruct{}
	SetTimestamps(ts, true)
	if ts.CreatedAt == 0 || ts.UpdatedAt == 0 {
		t.Errorf("Timestamps not set on create: CreatedAt=%d UpdatedAt=%d", ts.CreatedAt, ts.UpdatedAt)
	}

	time.Sleep(1000 * time.Millisecond)
	SetTimestamps(ts, false)
	if ts.UpdatedAt <= ts.CreatedAt {
		t.Errorf("UpdatedAt not updated on update: old=%d new=%d", ts.CreatedAt, ts.UpdatedAt)
	}
}

func TestSetTimestamps_Map(t *testing.T) {
	m := map[string]interface{}{}
	SetTimestamps(&m, true)
	if m["CreatedAt"] == nil || m["UpdatedAt"] == nil {
		t.Errorf("Timestamps not set in map: CreatedAt=%v UpdatedAt=%v", m["CreatedAt"], m["UpdatedAt"])
	}
}

func TestSetPKSKFromID_EmptyID(t *testing.T) {
	ts := &testStruct{}
	SetPKSKFromID(ts)
	if ts.PK != "" || ts.SK != "" {
		t.Errorf("Should not set PK/SK if ID is empty")
	}
}

func TestSetPKSKFromID_InvalidType(t *testing.T) {
	type badStruct struct{ ID int }
	bs := &badStruct{ID: 123}
	SetPKSKFromID(bs) // Should not panic
}

func TestSetIDFromPKSK_InvalidType(t *testing.T) {
	type badStruct struct {
		PK int
		SK int
	}
	bs := &badStruct{PK: 1, SK: 2}
	SetIDFromPKSK(bs) // Should not panic
}

func TestSetTimestamps_InvalidType(t *testing.T) {
	type badStruct struct {
		CreatedAt string
		UpdatedAt string
	}
	bs := &badStruct{}
	SetTimestamps(bs, true) // Should not panic
}

func TestGetFieldAndSetField(t *testing.T) {
	ts := &testStruct{}
	setField(reflect.ValueOf(ts).Elem(), "PK", "abc")
	f := getField(reflect.ValueOf(ts).Elem(), "PK", reflect.String)
	if !f.IsValid() || f.String() != "abc" {
		t.Errorf("getField/setField failed: got %v", f.String())
	}
}

// Struct só com PK/SK, sem ID
type onlyPKSKStruct struct {
	PK string
	SK string
}

func TestSetIDFromPKSK_OnlyPKSKStruct(t *testing.T) {
	s := &onlyPKSKStruct{PK: "p", SK: "s"}
	SetIDFromPKSK(s)
	v := reflect.ValueOf(s).Elem()
	idField := v.FieldByName("ID")
	if idField.IsValid() && idField.String() != "" {
		t.Errorf("Should not set ID in struct without ID field, got %v", idField.String())
	}
}

// Struct só com ID, sem PK/SK
type onlyIDStruct struct {
	ID string
}

func TestSetPKSKFromID_OnlyIDStruct(t *testing.T) {
	s := &onlyIDStruct{ID: "foo;bar"}
	SetPKSKFromID(s)
	// Não deve criar PK/SK em struct sem esses campos
	v := reflect.ValueOf(s).Elem()
	if _, ok := v.Type().FieldByName("PK"); ok {
		t.Errorf("PK should not exist in onlyIDStruct")
	}
	if _, ok := v.Type().FieldByName("SK"); ok {
		t.Errorf("SK should not exist in onlyIDStruct")
	}
}

func TestSetPKSKFromID_MapCreatesKeys(t *testing.T) {
	m := map[string]interface{}{"ID": "foo;bar"}
	SetPKSKFromID(&m)
	if m["PK"] != "foo" || m["SK"] != "bar" {
		t.Errorf("PK/SK not created in map: got PK=%v SK=%v", m["PK"], m["SK"])
	}
}

func TestSetPKSKFromID_IDSemSeparador(t *testing.T) {
	ts := &testStruct{ID: "apenaspk"}
	SetPKSKFromID(ts)
	if ts.PK != "apenaspk" || ts.SK != "" {
		t.Errorf("PK/SK not set correctly for ID without ';': got PK=%s, SK=%s", ts.PK, ts.SK)
	}
}

func TestSetPKSKFromID_IDComSeparador(t *testing.T) {
	ts := &testStruct{ID: "pkval;skval"}
	SetPKSKFromID(ts)
	if ts.PK != "pkval" || ts.SK != "skval" {
		t.Errorf("PK/SK not set correctly for ID with ';': got PK=%s, SK=%s", ts.PK, ts.SK)
	}
}

func TestSetIDFromPKSK_SomentePK(t *testing.T) {
	ts := &testStruct{PK: "apenaspk", SK: ""}
	SetIDFromPKSK(ts)
	if ts.ID != "apenaspk" {
		t.Errorf("ID not set correctly for only PK: got %s", ts.ID)
	}
}

func TestSetIDFromPKSK_PKandSK(t *testing.T) {
	ts := &testStruct{PK: "pkval", SK: "skval"}
	SetIDFromPKSK(ts)
	if ts.ID != "pkval;skval" {
		t.Errorf("ID not set correctly for PK and SK: got %s", ts.ID)
	}
}

func TestSetPKSKFromID_SomenteID(t *testing.T) {
	ts := &testStructWithId{ID: "pkonly"}
	SetPKSKFromID(ts)
	// Não deve criar PK/SK em struct sem esses campos
	v := reflect.ValueOf(ts).Elem()
	if _, ok := v.Type().FieldByName("PK"); ok {
		t.Errorf("PK should not exist in testStructWithId")
	}
	if _, ok := v.Type().FieldByName("SK"); ok {
		t.Errorf("SK should not exist in testStructWithId")
	}
}

func TestSetPKSKFromID_SomenteID_Map(t *testing.T) {
	m := map[string]interface{}{"ID": "pkonly"}
	SetPKSKFromID(&m)
	if m["PK"] != "pkonly" {
		t.Errorf("PK not created in map for ID without ';': got PK=%v", m["PK"])
	}
	if _, ok := m["SK"]; ok && m["SK"] != "" {
		t.Errorf("SK should be empty or not set for ID without ';': got SK=%v", m["SK"])
	}
}

func TestSetPKSKFromID_IDComSeparador_Map(t *testing.T) {
	m := map[string]interface{}{"ID": "pkval;skval"}
	SetPKSKFromID(&m)
	if m["PK"] != "pkval" || m["SK"] != "skval" {
		t.Errorf("PK/SK not created in map for ID with ';': got PK=%v SK=%v", m["PK"], m["SK"])
	}
}

func TestSetIDFromPKSK_SomentePKStruct(t *testing.T) {
	ts := &testStructWithPk{PK: "pkonly"}
	SetIDFromPKSK(ts)
	v := reflect.ValueOf(ts).Elem()
	idField := v.FieldByName("ID")
	if idField.IsValid() && idField.String() != "" {
		t.Errorf("Should not set ID in struct without ID field, got %v", idField.String())
	}
}

func TestSetIDFromPKSK_SomentePK_Map(t *testing.T) {
	m := map[string]interface{}{"PK": "pkonly"}
	SetIDFromPKSK(&m)
	if m["ID"] != "pkonly" {
		t.Errorf("ID not set correctly in map for only PK: got %v", m["ID"])
	}
}

func TestSetIDFromPKSK_PKSKStruct(t *testing.T) {
	ts := &testStructWithPkSk{PK: "pkval", SK: "skval"}
	SetIDFromPKSK(ts)
	v := reflect.ValueOf(ts).Elem()
	idField := v.FieldByName("ID")
	if idField.IsValid() && idField.String() != "" {
		t.Errorf("Should not set ID in struct without ID field, got %v", idField.String())
	}
}

func TestSetIDFromPKSK_PKSK_Map(t *testing.T) {
	m := map[string]interface{}{"PK": "pkval", "SK": "skval"}
	SetIDFromPKSK(&m)
	if m["ID"] != "pkval;skval" {
		t.Errorf("ID not set correctly in map for PK and SK: got %v", m["ID"])
	}
}
