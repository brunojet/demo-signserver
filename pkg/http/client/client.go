package http_client

import (
	"net/http"
	"time"
)

const (
	httpClientMinTries    = 1
	httpClientMaxTries    = 5
	httpClientMinInterval = 15
	httpClientMaxInterval = 60
)

type HttpClient struct {
	Adapter        HttpClientAdapter
	Tries          int
	Interval       int
	ShouldContinue func(statusCode int) bool // Corrigido nome
}

func NewHttpClient(adapter HttpClientAdapter) *HttpClient {
	if adapter == nil {
		panic("HttpClientAdapter cannot be nil")
	}
	return &HttpClient{
		Adapter:  adapter,
		Tries:    httpClientMinTries,
		Interval: httpClientMinInterval,
		ShouldContinue: func(statusCode int) bool {
			answer := false
			switch statusCode {
			case http.StatusOK, http.StatusCreated:
			case http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound, http.StatusConflict:
				break
			default:
				answer = true
			}

			return answer
		},
	}
}

func (h *HttpClient) SetShouldContinue(tries, interval int, shouldContinue func(statusCode int) bool) {
	if tries >= httpClientMinTries && tries <= httpClientMaxTries {
		h.Tries = tries
	}

	if interval >= httpClientMinInterval && interval <= httpClientMaxInterval {
		h.Interval = interval
	}

	if shouldContinue != nil {
		h.ShouldContinue = shouldContinue
	}
}

func (h *HttpClient) UploadFile(method HttpMethod, headers map[string]string, url string, path string) HttpClientResponse {
	for i := 0; i < h.Tries; i++ {
		response := h.Adapter.UploadFile(method, headers, url, path)
		if !h.ShouldContinue(response.StatusCode) {
			return response
		}
		time.Sleep(time.Duration(h.Interval) * time.Second)
	}
	return HttpClientResponse{Body: "", StatusCode: http.StatusGatewayTimeout}
}

func (h *HttpClient) DownloadFile(headers map[string]string, url string, dstPath string) HttpClientResponse {
	for i := 0; i < h.Tries; i++ {
		response := h.Adapter.DownloadFile(headers, url, dstPath)
		if !h.ShouldContinue(response.StatusCode) {
			return response
		}
		time.Sleep(time.Duration(h.Interval) * time.Second)
	}
	return HttpClientResponse{Body: "", StatusCode: http.StatusGatewayTimeout}
}
