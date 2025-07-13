package adapters

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"

	http_client "demo-signserver/pkg/http/client"
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

func (h *HttpClientAdapter) DownloadFile(headers map[string]string, url string, path string, response any) http_client.StatusCode {
	statusCode := http_client.StatusCode(500)

	req, err := newHttpRequest(string(http_client.HttpMethodGet), url, headers, nil)
	if err != nil {
		return statusCode
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return statusCode
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		out, err := os.Create(path)
		if err != nil {
			return statusCode
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return statusCode
		}
		return http_client.StatusCode(resp.StatusCode)
	default:
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, response)
		return http_client.StatusCode(resp.StatusCode)
	}
}

func (h *HttpClientAdapter) UploadFile(method http_client.HttpMethod, headers map[string]string, url string, path string, response any) http_client.StatusCode {
	statusCode := http_client.StatusCode(500)

	uploadFile, err := os.Open(path)
	if err != nil {
		return statusCode
	}
	defer uploadFile.Close()

	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)

	fileContentType := "binary/octet-stream"

	if headers != nil {
		if v, ok := headers["Content-Type"]; ok {
			fileContentType = v
			delete(headers, "Content-Type")
		}
	}

	go func() {
		var part io.Writer
		var err error
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition",
			`form-data; name="file"; filename="`+filepath.Base(path)+`"`)
		h.Set("Content-Type", fileContentType)
		part, err = mw.CreatePart(h)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		_, err = io.Copy(part, uploadFile)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		mw.Close()
		pw.Close()
	}()

	putHeaders := map[string]string{
		"Accept":       headers["Accept"],
		"Content-Type": mw.FormDataContentType(),
	}

	req, err := newHttpRequest(string(method), url, putHeaders, pr)
	if err != nil {
		return statusCode
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return statusCode
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if response != nil && len(body) > 0 {
		if err := json.Unmarshal(body, response); err != nil {
			return statusCode
		}
	}

	return http_client.StatusCode(resp.StatusCode)
}
