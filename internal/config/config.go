package config

import (
	"context"
	http_client "demo-signserver/pkg/http/client"
	http_client_adapters "demo-signserver/pkg/http/client/adapters"
	message_adapters "demo-signserver/pkg/message/adapters"
	"demo-signserver/pkg/observability"
	db_services "demo-signserver/pkg/repository/services"
	"demo-signserver/pkg/storage"
	storage_adapters "demo-signserver/pkg/storage/adapters"
	"fmt"
	"os"
	"sync"

	"demo-signserver/pkg/eventbus"
)

var (
	SignServerConfigInstance  *SignServerConfig
	SignServerMethodsInstance *SignServerMethods
	onceConfig                sync.Once
	onceMethods               sync.Once
)

type SignServerMethods struct {
	NewStorageService      func() storage.StorageAdapter
	NewDynamoDBService     func(tableName string, pkKey string, skKey string) *db_services.DynamoDBService
	NewMessageQueueAdapter func(unsignedDir, bucket string) message_adapters.MessageQueueAdapterInterface
	NewEventBus            func() *eventbus.EventBus
	NewHttpClient          func() http_client.HttpClientAdapter
}

type SignServerConfig struct {
	RequestTableName  string
	ProfileTableName  string
	StorageBucketName string
}

func OsGetenvPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s is not set", key))
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

		SignServerMethodsInstance = &SignServerMethods{
			NewMessageQueueAdapter: func(bucketName, unsignedDir string) message_adapters.MessageQueueAdapterInterface {
				// Replace with a real MetricsSink if available
				return message_adapters.NewLocalS3EventQueue(bucketName, unsignedDir, sink)
			},
			NewEventBus: func() *eventbus.EventBus {
				return eventbus.NewEventBusWithSink(sink)
			},
			NewStorageService: func() storage.StorageAdapter {
				if environment := os.Getenv("ENVIRONMENT"); environment == "local" {
					adapter := storage_adapters.NewLocalStorageService()
					return storage.NewStorageService(adapter)
				}

				panic("Storage service not implemented for non-local environments")
			},
			NewDynamoDBService: func(tableName string, pkKey string, skKey string) *db_services.DynamoDBService {
				return db_services.NewDynamoDBService(tableName, pkKey, skKey)
			},
			NewHttpClient: func() http_client.HttpClientAdapter {
				var adapter http_client.HttpClientAdapter
				if environment := os.Getenv("ENVIRONMENT"); environment == "local" {
					adapter = http_client_adapters.NewFileHttpClientAdapter()
				} else {
					adapter = http_client_adapters.NewHttpClientAdapter()
				}
				return http_client.NewHttpClient(adapter)
			},
		}
	})
	return SignServerMethodsInstance
}
