package adapters

import (
	"fmt"
	"io"
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

// GetBucketName implements StorageServiceInterface.
func (l *LocalStorageService) GetBucketName() string {
	return l.BasePath
}

// GeneratePresignedGetURL implements storage_services.StorageServiceInterface.
func (l *LocalStorageService) GeneratePresignedURL(httpMethod HttpMethod, key string, expires time.Duration) (string, error) {
	panic("unimplemented")
}

func NewLocalStorageService(basePath string) StorageServiceInterface {
	basePath = filepath.Join(os.TempDir(), basePath)
	return &LocalStorageService{BasePath: basePath}
}

// DownloadFile copia do "bucket" local para destino (ex: /tmp)
func (l *LocalStorageService) DownloadFile(key, dest string) error {
	srcPath := filepath.Join(l.BasePath, key)
	in, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo local: %w", err)
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório destino: %w", err)
	}
	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo destino: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// UploadFile copia do src para o "bucket" local
func (l *LocalStorageService) UploadFile(key, src string) error {
	destPath := filepath.Join(l.BasePath, key)
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório destino: %w", err)
	}
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo origem: %w", err)
	}
	defer in.Close()
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo destino: %w", err)
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
