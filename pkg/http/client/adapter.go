package http_client

type HttpMethod string

const (
	HttpMethodGet  HttpMethod = "GET"
	HttpMethodPut  HttpMethod = "PUT"
	HttpMethodPost HttpMethod = "POST"
)

type HttpClientResponse struct {
	Body       string
	StatusCode int
}

type HttpClientAdapter interface {
	UploadFile(method HttpMethod, headers map[string]string, url string, path string) HttpClientResponse
	DownloadFile(headers map[string]string, url string, path string) HttpClientResponse
}
