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

func MakeResourceName(resource string) string {
	return fmt.Sprintf("%s-%s-%s", OsGetenvPanic("PROJECT_NAME"), OsGetenvPanic("ENVIRONMENT"), OsGetenvPanic(resource))
}

func GetSignServerConfig() *SignServerConfig {
	once.Do(func() {
		SignServerConfigInstance = &SignServerConfig{
			RequestTableName:  MakeResourceName("SIGN_REQUEST_TABLE"),
			ProfileTableName:  MakeResourceName("SIGN_PROFILE_TABLE"),
			StorageBucketName: MakeResourceName("SIGN_STORAGE_BUCKET"),
		}
	})
	return SignServerConfigInstance
}
