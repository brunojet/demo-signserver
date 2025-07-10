package adapters

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var _ StorageServiceInterface = (*LocalStorageService)(nil)

// LocalStorageService simula operações de storage em disco local, compatível com S3ServiceInterface
// e métodos DownloadFile/UploadFile.
type LocalStorageService struct {
	BasePath string // diretório base simulando o bucket
	WorkPath string // caminho de trabalho opcional, se necessário
}

func NewLocalStorageService(basePath string) StorageServiceInterface {
	basePath = filepath.Join(os.TempDir(), basePath)
	workPath := filepath.Join(basePath, "work") // exemplo de caminho de trabalho
	if err := os.MkdirAll(filepath.Dir(workPath), 0755); err != nil {
		log.Fatalf("Erro ao criar diretório de trabalho: %v\n", err)
	}

	return &LocalStorageService{
		BasePath: basePath,
		WorkPath: workPath,
	}
}

// GetBucketName implements StorageServiceInterface.
func (l *LocalStorageService) GetBucketName() string {
	return l.BasePath
}

// GeneratePresignedGetURL implements storage_services.StorageServiceInterface.
func (l *LocalStorageService) GeneratePresignedURL(httpMethod HttpMethod, key string, expires time.Duration) (string, error) {
	panic("unimplemented")
}

func (l *LocalStorageService) DownloadFileFromS3(key string) error {
	srcPath := filepath.Join(l.BasePath, key)
	dstPath := filepath.Join(l.WorkPath, key)
	in, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo local: %w", err)
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		panic("erro ao criar diretório base: " + err.Error())
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo destino: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func (l *LocalStorageService) UploadToS3(key string) error {
	srcPath := filepath.Join(l.WorkPath, key)
	dstPath := filepath.Join(l.BasePath, key)

	in, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo origem: %w", err)
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		panic("erro ao criar diretório base: " + err.Error())
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo destino: %w", err)
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

// OpenWorkFile implements StorageServiceInterface.
func (l *LocalStorageService) OpenWorkFile(key string) (io.ReadWriteCloser, error) {
	workFilePath := filepath.Join(l.WorkPath, key)
	return os.Open(workFilePath)
}
