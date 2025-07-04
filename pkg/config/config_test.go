package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("PROJECT_NAME", "signserver")
	os.Setenv("ENVIRONMENT", "dev")
}

func TestNewConfigInstanceAndGetConfigInstance(t *testing.T) {
	name := "test1"
	cfg := NewConfigInstance(name)
	cfg2 := GetConfigInstance(name)
	if cfg2 == nil {
		t.Fatalf("Config instance '%s' not found", name)
	}
	if cfg != cfg2 {
		t.Errorf("Expected same config pointer, got different")
	}
}

func TestLoadDefaultConfig(t *testing.T) {
	cfg := LoadDefaultConfig()
	assert.NotNil(t, cfg, "Default config should not be nil")
	assert.NotEmpty(t, cfg.ProjectName, "Default config ProjectName should not be empty")
}

func TestSetAndGetResource(t *testing.T) {
	cfg := NewConfigInstance("resource-test")
	assert.NotNil(t, cfg, "Default config should not be nil")
	info := ResourceInfo{Name: "my-table", Type: "table", Parameters: map[string]string{"arn": "arn:aws:dynamodb:..."}}
	cfg.SetResource("Table1", info)
	res := cfg.GetResource("Table1")
	assert.NotNil(t, res, "Expected resource to be found after SetResource")
	// Como Parameters agora é any, precisa de type assertion
	assert.Equal(t, info.Name, res.Name, "ResourceInfo Name mismatch")
	assert.Equal(t, info.Type, res.Type, "ResourceInfo Type mismatch")
	params, ok := res.Parameters.(map[string]string)
	assert.True(t, ok, "Parameters should be map[string]string")
	assert.Equal(t, info.Parameters, params, "ResourceInfo Parameters mismatch")
}

func TestResourceNameBuilder(t *testing.T) {
	cfg := NewConfigInstance("builder-test")
	assert.NotNil(t, cfg, "Default config should not be nil")
	name := cfg.ResourceNameBuilder("my-resource")
	assert.NotEmpty(t, name, "ResourceNameBuilder should return a non-empty string")
}

func TestGetResourceNotFound(t *testing.T) {
	cfg := NewConfigInstance("notfound-test")
	assert.NotNil(t, cfg, "Default config should not be nil")
	cfg2 := cfg.GetResource("does-not-exist")
	assert.Nil(t, cfg2, "Expected GetResource to return nil for missing resource")
}
