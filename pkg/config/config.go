package config

import (
	"demo-signserver/pkg/observability"
	"fmt"
	"log"
	"os"
)

var (
	configInstances = make(map[string]*Config)
)

const (
	DefaultInstanceName = "default"
	defaultAwsRegion    = "us-east-1"
)

func LoadDefaultConfig() *Config {
	if _, exists := configInstances[DefaultInstanceName]; !exists {
		log.Printf("Default config instance '%s' not found, creating new instance", DefaultInstanceName)
		configInstances[DefaultInstanceName] = LoadConfig(DefaultInstanceName)
	} else {
		log.Printf("Using existing default config instance '%s'", DefaultInstanceName)
	}

	return configInstances[DefaultInstanceName]
}

type Config struct {
	ProjectName      string
	Environment      string
	AwsRegion        string
	DynamoDBEndpoint string
	Resources        map[string]ResourceInfo // Novo campo para recursos classificados
}

type ResourceInfo struct {
	Name       string
	Type       string // Ex: "table", "bucket", etc
	Parameters any    // Pode ser qualquer struct, map, etc
}

func getEnvEx(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func OsGetEnvPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}
func resourceNameBuilder(projectName, environment, resource string) string {
	return fmt.Sprintf("%s-%s-%s", projectName, environment, resource)
}

// Cria e armazena uma nova instância de configuração
func NewConfigInstance(name string) *Config {
	cfg := LoadConfig(name)
	configInstances[name] = cfg
	return cfg
}

// Recupera uma instância de configuração pelo nome
func GetConfigInstance(name string) *Config {
	if name == "" {
		name = DefaultInstanceName
	}
	cfg, ok := configInstances[name]

	if !ok {
		log.Fatalf("Config instance '%s' not found, creating new instance", name)
	}
	return cfg
}

func LoadConfig(instanceName string) *Config {
	projectName := OsGetEnvPanic("PROJECT_NAME")
	environment := OsGetEnvPanic("ENVIRONMENT")

	config := &Config{
		ProjectName:      projectName,
		Environment:      environment,
		AwsRegion:        getEnvEx("AWS_REGION", defaultAwsRegion),
		DynamoDBEndpoint: getEnvEx("DYNAMODB_ENDPOINT", ""),
		Resources:        make(map[string]ResourceInfo), // Inicializa vazio
	}

	observability.LogInfo("Config inicializada", map[string]interface{}{"config": config})

	return config
}

// Recupera um recurso pelo nome. Retorna o ResourceInfo e true se existir, ou false se não encontrado.
func (c *Config) GetResource(name string) *ResourceInfo {
	res, ok := c.Resources[name]
	if !ok {
		log.Printf("Resource '%s' not found in config instance '%s'", name, c.ProjectName)
		return nil
	}
	return &res
}

// Adiciona ou atualiza um recurso no map Resources
func (c *Config) SetResource(name string, info ResourceInfo) {
	if c.Resources == nil {
		c.Resources = make(map[string]ResourceInfo)
	}
	c.Resources[name] = info
}

// Atualiza ou adiciona parâmetros a um recurso existente
func (c *Config) SetResourceParameters(name string, params any) bool {
	res, ok := c.Resources[name]
	if !ok {
		return false
	}
	res.Parameters = params
	c.Resources[name] = res
	return true
}

func (c *Config) ResourceNameBuilder(resource string) string {
	return resourceNameBuilder(c.ProjectName, c.Environment, resource)
}
