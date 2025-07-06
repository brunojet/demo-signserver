package config

import (
	"context"
	db_services "demo-signserver/pkg/repository/services"
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

func SetupLocalEnvironment() {
	if environment := os.Getenv("ENVIRONMENT"); environment == "local" {
		db := db_services.NewDB()
		config := GetSignServerConfig()
		db.CreateTable(context.Background(), config.ProfileTableName, db_services.SORT_KEY)
		db.CreateTable(context.Background(), config.RequestTableName, db_services.NO_KEY)
		fmt.Println("Local environment setup completed.")
	}
}
