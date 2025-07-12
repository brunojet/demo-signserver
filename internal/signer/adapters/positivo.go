package adapters

import (
	"demo-signserver/internal/config"
	"demo-signserver/internal/repository/domain"
	"demo-signserver/internal/repository/repositories"
	http_client "demo-signserver/pkg/http/client"
	"fmt"
	"net/http"
)

type ExternalSigner interface {
	StartSign(contentType string, uploadFilePath string) error
	WaitSignature(ID string, downloadFilePath string) error
}

type PositivoSigner struct {
	HttpClientAdapter http_client.HttpClientAdapter
	SignerProfile     *domain.SignerProfile
}

func NewPositivoSigner(profileID string) (*PositivoSigner, error) {
	methods := config.GetSignServerMethods()
	httpClientAdapter := methods.NewHttpClient()
	repo := repositories.NewProfileRepository()
	profile, err := repo.GetProfileByID(profileID)
	if err != nil {
		return nil, err
	}

	return &PositivoSigner{
		SignerProfile:     profile,
		HttpClientAdapter: httpClientAdapter,
	}, nil
}

func (s *PositivoSigner) StartSign(srcPath string) (string, error) {
	endpoint := s.SignerProfile.Upload
	httpClient := http_client.NewHttpClient(s.HttpClientAdapter)
	httpClient.SetShouldContinue(endpoint.Tries, endpoint.Interval, nil)

	headers := map[string]string{
		"Content-Type": *s.SignerProfile.ContentType,
	}

	response := httpClient.UploadFile(http_client.HttpMethodGet, headers, endpoint.URL, srcPath)

	if response.StatusCode != http.StatusOK {
		return response.Body, fmt.Errorf("erro ao enviar arquivo: (status: %d)", response.StatusCode)
	}

	return response.Body, nil
}

func (s *PositivoSigner) WaitSignature(ID string, dstPath string) error {
	endpoint := s.SignerProfile.Download
	httpClient := http_client.NewHttpClient(s.HttpClientAdapter)
	httpClient.SetShouldContinue(endpoint.Tries, endpoint.Interval, nil)

	headers := map[string]string{
		"Accept": *s.SignerProfile.ContentType,
	}

	response := httpClient.DownloadFile(headers, fmt.Sprintf("%s/%s", endpoint.URL, ID), dstPath)

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("erro ao enviar arquivo: (status: %d)", response.StatusCode)
	}

	return nil
}
