package adapters

import (
	"crypto/md5"
	"demo-signserver/pkg/file"
	"demo-signserver/pkg/storage"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

var _ storage.StorageAdapter = (*LocalStorageService)(nil)

type LocalStorageService struct {
	WorkPath string
}

func NewLocalStorageService(localStorage string) storage.StorageAdapter {
	workPath := filepath.Join(localStorage, "local_storage_service")
	file.CreatePathIfNotExists(workPath)

	return &LocalStorageService{
		WorkPath: workPath,
	}
}

func (l *LocalStorageService) GeneratePresignedURL(httpMethod storage.HttpMethod, fileInfo *storage.FileInfo, expires time.Duration) (string, error) {
	return path.Join("https://example.com", filepath.ToSlash(fileInfo.FilePath)), nil
}

func (l *LocalStorageService) DownloadFileFromS3(fileInfo *storage.FileInfo) error {
	srcPath := filepath.Join(fileInfo.StoragePath, fileInfo.FilePath)
	dstPath := filepath.Join(l.WorkPath, fileInfo.FilePath)

	in, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo local: %w", err)
	}
	defer in.Close()

	out, err := file.CreateFilePath(dstPath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo destino: %w", err)
	}
	defer out.Close()

	hash := md5.New()
	mw := io.MultiWriter(out, hash)
	buf := make([]byte, 1024*1024) // 1MB buffer
	size, err := io.CopyBuffer(mw, in, buf)
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

	out, err := file.CreateFilePath(dstPath)

	if err != nil {
		return fmt.Errorf("erro ao criar arquivo destino: %w", err)
	}
	defer out.Close()
	buf := make([]byte, 1024*1024) // 1MB buffer
	_, err = io.CopyBuffer(out, in, buf)
	return err
}

// GetWorkFilePath implements StorageServiceInterface.
func (l *LocalStorageService) GetWorkFilePath(fileInfo *storage.FileInfo) string {
	return filepath.Join(l.WorkPath, fileInfo.FilePath)
}
