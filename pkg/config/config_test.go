package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigInstanceAndGetConfigInstance(t *testing.T) {
	name := "test1"
	cfg := NewConfigInstance(name)
	cfg2, ok := GetConfigInstance(name)
	if !ok {
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
	res, ok := cfg.GetResource("Table1")
	assert.True(t, ok, "Expected resource to be found after SetResource")
	assert.Equal(t, info, res, "ResourceInfo mismatch")
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
	_, ok := cfg.GetResource("does-not-exist")
	assert.False(t, ok, "Expected GetResource to return false for missing resource")
}
