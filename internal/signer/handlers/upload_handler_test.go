package handlers

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	"demo-signserver/internal/config"
	"demo-signserver/internal/repository/domain"
	"demo-signserver/internal/repository/repositories"
	storages "demo-signserver/internal/storage"
	"demo-signserver/pkg/eventbus"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	bus *eventbus.EventBus
)

func DummyHandler(bus *eventbus.EventBus) eventbus.Handler {
	return func(ctx context.Context, event any) {
		// Dummy handler for testing purposes
	}
}

func TestMain(m *testing.M) {
	os.Setenv("PROJECT_NAME", "signer-upload")
	os.Setenv("ENVIRONMENT", "local")
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8001")
	os.Setenv("SIGN_REQUEST_TABLE", "request")
	os.Setenv("SIGN_PROFILE_TABLE", "profile")
	os.Setenv("SIGN_STORAGE_BUCKET", "storage")
	os.Setenv("AWS_ACCESS_KEY_ID", "fake")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "fake")
	basePath := filepath.Join(os.TempDir(), config.MakeResourceName("SIGN_STORAGE_BUCKET"))
	if _, err := os.Stat(basePath); err == nil {
		os.RemoveAll(basePath)
	}
	if err := os.MkdirAll(basePath, 0755); err != nil {
		panic("erro ao criar diretório base: " + err.Error())
	}
	defer os.RemoveAll(basePath)
	bus = eventbus.NewEventBus()
	if err := bus.Register("upload_received", UploadReceivedHandler(bus), 1, 1); err != nil {
		log.Fatalf("Erro ao registrar handler: %v", err)
	}
	if err := bus.Register("sign_process", DummyHandler(bus), 1, 1); err != nil {
		log.Fatalf("Erro ao registrar handler: %v", err)
	}
	exitCode := m.Run()
	bus.Stop()
	os.Exit(exitCode)
}

func TestUploadReceivedHandler_Success(t *testing.T) {
	storage := storages.NewStorageService()
	bucket := storage.GetBucketName()
	key := filepath.Join("unsigned", uuid.New().String())
	sha := "sha256dummy"
	size := int64(11)

	filePath := filepath.Join(storage.GetBucketName(), key)

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		panic("erro ao criar diretório base: " + err.Error())
	}

	err := os.WriteFile(filePath, []byte("hello world"), 0644)
	assert.NoError(t, err, "Erro ao escrever arquivo de teste")

	// Setup: cria request no repositório fake
	repo := repositories.NewRequestRepository()
	request := &domain.SignRequest{}
	request.SetID(filepath.Base(key))
	request.SetSignerStatus(domain.SignerStatusCreated, nil)
	request.SetUnsignedBucketInfo(bucket, key)
	err = repo.CreateRequest(request)
	assert.NoError(t, err, "Erro ao criar request no repositório")
	handler := UploadReceivedHandler(bus)
	handler(context.Background(), UploadEvent{
		Bucket: bucket,
		Key:    key,
		ETag:   sha,
		Size:   size,
	})
}
