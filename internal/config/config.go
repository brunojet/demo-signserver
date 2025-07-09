package config

import (
	"context"
	message_adapters "demo-signserver/pkg/message/adapters"
	"demo-signserver/pkg/observability"
	db_services "demo-signserver/pkg/repository/services"
	storage_adapters "demo-signserver/pkg/storage/adapters"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	SignServerConfigInstance  *SignServerConfig
	SignServerMethodsInstance *SignServerMethods
	onceConfig                sync.Once
	onceMethods               sync.Once
)

type SignServerMethods struct {
	NewStorageService   func(bucketName string) storage_adapters.StorageServiceInterface
	NewDynamoDBService  func(tableName string, pkKey string, skKey string) *db_services.DynamoDBService
	MessageQueueAdapter func(unsignedDir, bucket string) message_adapters.MessageQueueAdapterInterface
}

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
	onceConfig.Do(func() {
		requestTable := MakeResourceName("SIGN_REQUEST_TABLE")
		profileTable := MakeResourceName("SIGN_PROFILE_TABLE")
		storageBucket := MakeResourceName("SIGN_STORAGE_BUCKET")

		SignServerConfigInstance = &SignServerConfig{
			ProfileTableName:  profileTable,
			RequestTableName:  requestTable,
			StorageBucketName: storageBucket,
		}

		if environment := os.Getenv("ENVIRONMENT"); environment == "local" {
			db := db_services.NewDB()
			db.CreateTable(context.Background(), profileTable, db_services.SORT_KEY)
			db.CreateTable(context.Background(), requestTable, db_services.NO_KEY)
		}
	})
	return SignServerConfigInstance
}

func GetSignServerMethods() *SignServerMethods {
	onceMethods.Do(func() {
		sink := observability.NewAccumulatorSink()
		metricsService := observability.NewMetricsService(sink)

		SignServerMethodsInstance = &SignServerMethods{
			MessageQueueAdapter: func(bucketName, unsignedDir string) message_adapters.MessageQueueAdapterInterface {
				// Replace with a real MetricsSink if available
				return message_adapters.NewLocalS3EventQueue(bucketName, unsignedDir, metricsService)
			},
			NewStorageService: func(bucketName string) storage_adapters.StorageServiceInterface {
				if environment := os.Getenv("ENVIRONMENT"); environment == "local" {
					return storage_adapters.NewLocalStorageService(bucketName)
				} else {
					return storage_adapters.NewS3Service(bucketName)
				}
			},
			NewDynamoDBService: func(tableName string, pkKey string, skKey string) *db_services.DynamoDBService {
				return db_services.NewDynamoDBService(tableName, pkKey, skKey)
			},
		}
	})
	return SignServerMethodsInstance
}
