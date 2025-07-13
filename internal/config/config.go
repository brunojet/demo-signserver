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
	"log"
	"os"
	"path/filepath"
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
	NewStorageService  func(ctx context.Context) storage.StorageAdapter
	NewDynamoDBService func(tableName string, pkKey string, skKey string) *db_services.DynamoDBService
	NewMessageQueue    func(ctx context.Context) message_adapters.MessageQueueAdapter
	NewEventBus        func() *eventbus.EventBus
	NewHttpClient      func(ctx context.Context) *http_client.HttpClient
}

type SignServerConfig struct {
	RequestTableName  string
	ProfileTableName  string
	StorageBucketName string
	LocalStorage      *string
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

		var localStorage *string

		if OsGetenvPanic("ENVIRONMENT") == "local" {
			localPath := filepath.Join(os.TempDir(), "demo_sign_server")
			storageBucket = filepath.Join(localPath, "sign_storage_bucket")
			localStorage = &localPath
			log.Printf("Using local storage path: %s\n", *localStorage)
		}

		SignServerConfigInstance = &SignServerConfig{
			ProfileTableName:  profileTable,
			RequestTableName:  requestTable,
			StorageBucketName: storageBucket,
			LocalStorage:      localStorage,
		}

		if OsGetenvPanic("ENVIRONMENT") == "local" {
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
		config := GetSignServerConfig()

		SignServerMethodsInstance = &SignServerMethods{
			NewMessageQueue: func(ctx context.Context) message_adapters.MessageQueueAdapter {
				ctx = observability.ContextWithSink(ctx, sink)
				if config.LocalStorage != nil {
					return message_adapters.NewLocalS3EventQueue(ctx, config.StorageBucketName, "unsigned")
				}
				panic("Message queue adapter not implemented for non-local environments")
			},
			NewEventBus: func() *eventbus.EventBus {
				return eventbus.NewEventBusWithSink(sink)
			},
			NewStorageService: func(ctx context.Context) storage.StorageAdapter {
				ctx = observability.ContextWithSink(ctx, sink)
				var adapter storage.StorageAdapter
				if config.LocalStorage != nil {
					adapter = storage_adapters.NewLocalStorageService(*config.LocalStorage)
				} else {
					panic("Storage adapter not implemented for non-local environments")
				}
				return storage.NewStorageService(ctx, adapter)
			},
			NewDynamoDBService: func(tableName string, pkKey string, skKey string) *db_services.DynamoDBService {
				return db_services.NewDynamoDBService(tableName, pkKey, skKey)
			},
			NewHttpClient: func(ctx context.Context) *http_client.HttpClient {
				ctx = observability.ContextWithSink(ctx, sink)
				var adapter http_client.HttpClientAdapter
				if config.LocalStorage != nil {
					adapter = http_client_adapters.NewFileHttpClientAdapter(*config.LocalStorage)
				} else {
					adapter = http_client_adapters.NewHttpClientAdapter()
				}
				return http_client.NewHttpClient(ctx, adapter)
			},
		}
	})
	return SignServerMethodsInstance
}
