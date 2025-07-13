package adapters

import (
	"context"
	"demo-signserver/internal/config"
	"demo-signserver/internal/repository/domain"
	http_client "demo-signserver/pkg/http/client"
	"fmt"
	"net/http"
)

type ExternalSigner interface {
	StartSign(contentType string, uploadFilePath string) error
	WaitSignature(ID string, downloadFilePath string) error
}

type PositivoSigner struct {
	HttpClient    *http_client.HttpClient
	SignerProfile *domain.SignerProfile
}

func NewPositivoSigner(ctx context.Context, profile *domain.SignerProfile) (ExternalSignerAdapter, error) {
	methods := config.GetSignServerMethods()
	httpClient := methods.NewHttpClient(ctx)

	return &PositivoSigner{
		SignerProfile: profile,
		HttpClient:    httpClient,
	}, nil
}

type PositivoSignerResponse struct {
	ID string `json:"id"`
}

func (s *PositivoSigner) StartSign(srcPath string) (string, error) {
	endpoint := s.SignerProfile.Upload
	s.HttpClient.SetShouldContinue(endpoint.Tries, endpoint.Interval, nil)

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "vnd.android.package-archive",
	}

	response := &PositivoSignerResponse{}

	statusCode := s.HttpClient.UploadFile(http_client.HttpMethodGet, headers, endpoint.URL, srcPath, response)

	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		return "", fmt.Errorf("erro ao enviar arquivo: (status: %d)", statusCode)
	}

	return response.ID, nil
}

func (s *PositivoSigner) WaitSignature(ID string, dstPath string) error {
	endpoint := s.SignerProfile.Download
	s.HttpClient.SetShouldContinue(endpoint.Tries, endpoint.Interval, nil)

	headers := map[string]string{
		"Accept": "vnd.android.package-archive",
	}

	response := &PositivoSignerResponse{}

	statusCode := s.HttpClient.DownloadFile(headers, fmt.Sprintf("%s/%s", endpoint.URL, ID), dstPath, response)

	if statusCode != http.StatusOK {
		return fmt.Errorf("erro ao receber arquivo: (status: %d)", statusCode)
	}

	return nil
}
