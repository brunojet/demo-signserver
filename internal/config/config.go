package config

import (
	"demo-signserver/pkg/config"
	"demo-signserver/pkg/observability"
)

func ConfigInit() {
	cfg := config.LoadDefaultConfig()

	profile_table := &config.ResourceInfo{
		Name: cfg.ResourceNameBuilder(config.OsGetEnvPanic("SIGN_PROFILE_TABLE")),
		Type: "dynamodb",
	}

	cfg.SetResource("profile_table", *profile_table)

	request_table := &config.ResourceInfo{
		Name: cfg.ResourceNameBuilder(config.OsGetEnvPanic("SIGN_REQUEST_TABLE")),
		Type: "dynamodb",
	}

	cfg.SetResource("request_table", *request_table)

	storage := &config.ResourceInfo{
		Name: cfg.ResourceNameBuilder(config.OsGetEnvPanic("SIGN_STORAGE_BUCKET")),
		Type: "s3",
	}

	cfg.SetResource("storage", *storage)

	observability.LogInfo("ConfigResource", map[string]interface{}{
		"name":       cfg.ProjectName,
		"info":       cfg.Environment,
		"aws_region": cfg.AwsRegion,
	})
}
