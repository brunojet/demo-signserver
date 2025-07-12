package adapters

import (
	"crypto/md5"
	"demo-signserver/pkg/storage"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var _ storage.StorageAdapter = (*LocalStorageService)(nil)

type LocalStorageService struct {
	WorkPath string
}

func NewLocalStorageService() storage.StorageAdapter {
	workPath := filepath.Join(os.TempDir(), "local_storage_service")
	if err := os.MkdirAll(filepath.Dir(workPath), 0755); err != nil {
		log.Fatalf("Erro ao criar diretório de trabalho: %v\n", err)
	}

	return &LocalStorageService{
		WorkPath: workPath,
	}
}

func (l *LocalStorageService) GeneratePresignedURL(httpMethod storage.HttpMethod, fileInfo *storage.FileInfo, expires time.Duration) (string, error) {
	return "https://example.com/" + fileInfo.FilePath, nil
}

func (l *LocalStorageService) DownloadFileFromS3(fileInfo *storage.FileInfo) error {
	srcPath := filepath.Join(fileInfo.StoragePath, fileInfo.FilePath)
	dstPath := filepath.Join(l.WorkPath, fileInfo.FilePath)
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

	hash := md5.New()
	mw := io.MultiWriter(out, hash)
	size, err := io.Copy(mw, in)
	if err != nil {
		return fmt.Errorf("erro ao copiar arquivo: %w", err)
	}

	fileInfo.Hash = fmt.Sprintf("%x", hash.Sum(nil))
	fileInfo.Size = size
	fileInfo.State = storage.FileStateReady

	return nil
}

func (l *LocalStorageService) UploadToS3(fileInfo *storage.FileInfo) error {
	srcPath := filepath.Join(l.WorkPath, fileInfo.FilePath)
	dstPath := filepath.Join(fileInfo.StoragePath, fileInfo.FilePath)

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

// GetWorkFilePath implements StorageServiceInterface.
func (l *LocalStorageService) GetWorkFilePath(fileInfo *storage.FileInfo) string {
	return filepath.Join(l.WorkPath, fileInfo.FilePath)
}
