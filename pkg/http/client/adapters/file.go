package adapters

import (
	http_client "demo-signserver/pkg/http/client"
	"io"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

var _ http_client.HttpClientAdapter = (*FileHttpClientAdapter)(nil)

type FileHttpClientAdapter struct {
	WorkPath string
	Files    map[string]fileMeta
}

type fileMeta struct {
	Path      string
	ReadyTime time.Time
}

func NewFileHttpClientAdapter() http_client.HttpClientAdapter {
	workPath := filepath.Join(os.TempDir(), "http_file_client")

	if _, err := os.Stat(workPath); os.IsNotExist(err) {
		if err := os.MkdirAll(workPath, 0755); err != nil {
			panic("Failed to create work path: " + err.Error())
		}
	}

	return &FileHttpClientAdapter{
		WorkPath: workPath,
		Files:    make(map[string]fileMeta),
	}
}

func (h *FileHttpClientAdapter) UploadFile(_ http_client.HttpMethod, _ map[string]string, _ string, srcPath string) http_client.HttpClientResponse {
	response := http_client.HttpClientResponse{Body: "", StatusCode: 500}
	src, err := os.Open(srcPath)
	if err != nil {
		return response
	}
	defer src.Close()

	dstPath := filepath.Join(h.WorkPath, uuid.New().String())
	dst, err := os.Create(dstPath)
	if err != nil {
		return response
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return response
	}

	id := path.Base(dstPath)
	delay := rand.Intn(31) + 15
	meta := fileMeta{Path: dstPath, ReadyTime: time.Now().Add(time.Duration(delay) * time.Second)}
	h.Files[id] = meta

	return response
}

func (h *FileHttpClientAdapter) DownloadFile(_ map[string]string, downloadURL string, dstPath string) http_client.HttpClientResponse {
	response := http_client.HttpClientResponse{Body: "", StatusCode: 500}
	id := path.Base(downloadURL)
	meta, ok := h.Files[id]
	if !ok {
		response.StatusCode = 404
		return response
	}

	src, err := os.Open(meta.Path)
	if err != nil {
		response.StatusCode = 404
		return response
	}
	defer src.Close()

	if time.Now().Before(meta.ReadyTime) {
		response.StatusCode = 202
		return response
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return response
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return response
	}

	response.StatusCode = 200
	return response
}
