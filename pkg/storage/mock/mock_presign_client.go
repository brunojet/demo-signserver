// MockPresignClient é um mock do PresignClient do S3 para uso em testes.
// Permite simular respostas customizadas para PresignPutObject.
// Exemplo de uso:
//
//	mock := &MockPresignClient{URL: "https://mock-url"}
package storage_mock

import (
	"context"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MockPresignClient struct {
	URL string
	Err error
}

// PresignPutObject implementa a interface do PresignClient do S3, usando o comportamento customizado se definido.
func (m *MockPresignClient) PresignPutObject(_ context.Context, _ *s3.PutObjectInput, _ ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return &v4.PresignedHTTPRequest{URL: m.URL}, nil
}
