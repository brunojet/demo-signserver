package http_client

import (
	"context"
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
	ctx            context.Context
	obs            *ObservabilityHttpClientMiddleware
	adapter        HttpClientAdapter
	tries          int
	interval       int
	shouldContinue func(statusCode StatusCode) bool
}

func NewHttpClient(ctx context.Context, adapter HttpClientAdapter) *HttpClient {
	if adapter == nil {
		panic("HttpClientAdapter cannot be nil")
	}
	return &HttpClient{
		ctx:      ctx,
		obs:      NewObservabilityHttpClientMiddleware(ctx),
		adapter:  adapter,
		tries:    httpClientMinTries,
		interval: httpClientMinInterval,
		shouldContinue: func(statusCode StatusCode) bool {
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

func (h *HttpClient) SetShouldContinue(tries, interval int, shouldContinue func(statusCode StatusCode) bool) {
	if tries >= httpClientMinTries && tries <= httpClientMaxTries {
		h.tries = tries
	}

	if interval >= httpClientMinInterval && interval <= httpClientMaxInterval {
		h.interval = interval
	}

	if shouldContinue != nil {
		h.shouldContinue = shouldContinue
	}
}

func (h *HttpClient) UploadFile(method HttpMethod, headers map[string]string, url string, path string, response any) StatusCode {
	for i := 0; i < h.tries; i++ {
		statusCode := h.obs.do("UploadFile", func() StatusCode {
			return h.adapter.UploadFile(method, headers, url, path, response)
		})
		if !h.shouldContinue(statusCode) {
			return statusCode
		}
		time.Sleep(time.Duration(h.interval) * time.Second)
	}
	return StatusCode(http.StatusGatewayTimeout)
}

func (h *HttpClient) DownloadFile(headers map[string]string, url string, dstPath string, response any) StatusCode {
	for i := 0; i < h.tries; i++ {
		statusCode := h.obs.do("DownloadFile", func() StatusCode {
			return h.adapter.DownloadFile(headers, url, dstPath, response)
		})

		if !h.shouldContinue(statusCode) {
			return statusCode
		}
		time.Sleep(time.Duration(h.interval) * time.Second)
	}
	return StatusCode(http.StatusGatewayTimeout)
}
