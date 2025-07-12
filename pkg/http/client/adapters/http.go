package adapters

import (
	http_client "demo-signserver/pkg/http/client"
	"io"
	"net/http"
	"os"
)

var _ http_client.HttpClientAdapter = (*HttpClientAdapter)(nil)

type HttpClientAdapter struct {
}

func NewHttpClientAdapter() http_client.HttpClientAdapter {
	return &HttpClientAdapter{}
}

func newHttpRequest(method, url string, headers map[string]string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req, nil
}

func (h *HttpClientAdapter) DownloadFile(headers map[string]string, url string, path string) http_client.HttpClientResponse {
	response := http_client.HttpClientResponse{Body: "", StatusCode: 500}

	req, err := newHttpRequest(string(http_client.HttpMethodGet), url, headers, nil)
	if err != nil {
		return response
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return response
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		out, err := os.Create(path)
		if err != nil {
			return response
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return response
		}
		response.StatusCode = 200
		return response
	default:
		body, _ := io.ReadAll(resp.Body)
		response.Body = string(body)
		response.StatusCode = resp.StatusCode
		return response
	}
}

func (h *HttpClientAdapter) UploadFile(method http_client.HttpMethod, headers map[string]string, url string, path string) http_client.HttpClientResponse {
	response := http_client.HttpClientResponse{Body: "", StatusCode: 500}

	uploadFile, err := os.Open(path)
	if err != nil {
		return response
	}
	defer uploadFile.Close()

	req, err := newHttpRequest(string(method), url, headers, uploadFile)
	if err != nil {
		return response
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return response
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	response.Body = string(body)
	response.StatusCode = resp.StatusCode

	return response
}
