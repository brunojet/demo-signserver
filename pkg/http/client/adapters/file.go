package adapters

import (
	"demo-signserver/pkg/file"
	http_client "demo-signserver/pkg/http/client"
	"encoding/json"
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

func NewFileHttpClientAdapter(localStorage string) http_client.HttpClientAdapter {
	workPath := filepath.Join(localStorage, "http_file_client")
	file.CreatePathIfNotExists(workPath)

	return &FileHttpClientAdapter{
		WorkPath: workPath,
		Files:    make(map[string]fileMeta),
	}
}

func (h *FileHttpClientAdapter) UploadFile(_ http_client.HttpMethod, _ map[string]string, _ string, srcPath string, response any) http_client.StatusCode {
	statusCode := http_client.StatusCode(500)

	src, err := os.Open(srcPath)
	if err != nil {
		return statusCode
	}
	defer src.Close()

	dstPath := filepath.Join(h.WorkPath, uuid.New().String())

	dst, err := file.CreateFilePath(dstPath)

	if err != nil {
		return statusCode
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return statusCode
	}

	id := filepath.Base(dstPath)
	delay := rand.Intn(30) + 15
	meta := fileMeta{Path: dstPath, ReadyTime: time.Now().Add(time.Duration(delay) * time.Second)}
	h.Files[id] = meta

	bodyMap := map[string]string{"id": id}
	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return statusCode
	}

	if response != nil {
		_ = json.Unmarshal(bodyBytes, response)
	}

	return http_client.StatusCode(200)
}

func (h *FileHttpClientAdapter) DownloadFile(_ map[string]string, downloadURL string, dstPath string, response any) http_client.StatusCode {
	statusCode := http_client.StatusCode(500)
	id := path.Base(downloadURL)
	meta, ok := h.Files[id]
	if !ok {
		return statusCode
	}

	src, err := os.Open(meta.Path)
	if err != nil {
		return http_client.StatusCode(404)
	}
	defer src.Close()

	if time.Now().Before(meta.ReadyTime) {
		return http_client.StatusCode(202)
	}

	dst, err := file.CreateFilePath(dstPath)
	if err != nil {
		return statusCode
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return statusCode
	}

	return http_client.StatusCode(200)
}
