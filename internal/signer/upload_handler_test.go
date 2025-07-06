package signer

import (
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
	exitCode := m.Run()
	bus.Stop()
	os.Exit(exitCode)
}

func TestUploadReceivedHandler_Success(t *testing.T) {
	storage := storages.NewStorageService()
	bucket := storage.GetBucketName()
	key := uuid.New().String()
	sha := "sha256dummy"
	size := int64(11)

	filePath := filepath.Join(storage.GetBucketName(), key)
	os.WriteFile(filePath, []byte("hello world"), 0644)

	// Setup: cria request no repositório fake
	repo := repositories.NewRequestRepository()
	request := &domain.SignRequest{}
	request.SetID(key)
	request.SetUnsignedBucketInfo(bucket, key)
	repo.CreateRequest(request)

	// Setup: EventBus fake
	bus.Publish("upload_received", UploadEvent{
		Bucket: bucket,
		Key:    key,
		SHA256: sha,
		Size:   size,
	})
	err := bus.WaitForHandlers("upload_received")
	assert.NoError(t, err, "Erro ao esperar handlers")
	savedFilePath := filepath.Join(storage.GetBucketName(), key)
	info, err := os.Stat(savedFilePath)
	assert.NoError(t, err, "Arquivo não foi salvo no bucket")
	assert.False(t, info.IsDir(), "O caminho salvo não é um diretório")
	assert.Equal(t, size, info.Size(), "Tamanho do arquivo salvo está incorreto")
}
