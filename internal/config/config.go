package config

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	SignServerConfigInstance *SignServerConfig
	once                     sync.Once
)

type SignServerConfig struct {
	RequestTableName  string
	ProfileTableName  string
	StorageBucketName string
}

func OsGetenvPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

func makeResourceName(projectName, environment, resource string) string {
	return fmt.Sprintf("%s-%s-%s", projectName, environment, OsGetenvPanic(resource))
}

func GetSignServerConfig() *SignServerConfig {
	once.Do(func() {
		profileName := OsGetenvPanic("PROJECT_NAME")
		environment := OsGetenvPanic("ENVIRONMENT")
		SignServerConfigInstance = &SignServerConfig{
			RequestTableName:  makeResourceName(profileName, environment, "SIGN_REQUEST_TABLE"),
			ProfileTableName:  makeResourceName(profileName, environment, "SIGN_PROFILE_TABLE"),
			StorageBucketName: makeResourceName(profileName, environment, "SIGN_STORAGE_BUCKET"),
		}
	})
	return SignServerConfigInstance
}
