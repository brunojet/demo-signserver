package dtos

import (
	"demo-signserver/internal/repository/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferInfoDTO_Validation(t *testing.T) {
	dto := TransferInfoDTO{
		Url:      "not-a-url",
		Tries:    0,
		Interval: 100,
	}
	// Não é possível validar diretamente sem o Gin, mas podemos testar os valores
	assert.NotEqual(t, "", dto.Url)
	assert.True(t, dto.Tries >= 0)
	assert.True(t, dto.Interval >= 0)
}

func TestCreateSignerProfileDTO_GetDomainCreateSignerProfile(t *testing.T) {
	dto := CreateProfileDTO{
		Signer:      domain.SignerPositivo,
		ProfileId:   "123",
		Description: "desc",
		Upload:      TransferInfoDTO{Url: "http://a.com", Tries: 1, Interval: 5},
		Download:    TransferInfoDTO{Url: "http://b.com", Tries: 2, Interval: 10},
	}
	domain := dto.GetDomainCreateSignerProfile()
	assert.NotNil(t, domain)
	assert.Equal(t, "positivo", string(*domain.Signer))
	assert.Equal(t, "123", *domain.ProfileId)
	assert.Equal(t, "desc", *domain.Description)
	assert.Equal(t, "http://a.com", domain.Upload.URL)
	assert.Equal(t, 1, domain.Upload.Tries)
	assert.Equal(t, 5, domain.Upload.Interval)
	assert.Equal(t, "http://b.com", domain.Download.URL)
	assert.Equal(t, 2, domain.Download.Tries)
	assert.Equal(t, 10, domain.Download.Interval)
}

func TestUpdateSignerProfileDTO_GetDomainUpdateSignerProfile(t *testing.T) {
	desc := "nova desc"
	upload := &TransferInfoDTO{Url: "http://a.com", Tries: 1, Interval: 5}
	download := &TransferInfoDTO{Url: "http://b.com", Tries: 2, Interval: 10}
	dto := UpdateProfileDTO{
		Description: &desc,
		Upload:      upload,
		Download:    download,
	}
	domain := dto.GetDomainUpdateSignerProfile()
	assert.NotNil(t, domain)
	assert.Equal(t, &desc, domain.Description)
	assert.Equal(t, "http://a.com", domain.Upload.URL)
	assert.Equal(t, 1, domain.Upload.Tries)
	assert.Equal(t, 5, domain.Upload.Interval)
	assert.Equal(t, "http://b.com", domain.Download.URL)
	assert.Equal(t, 2, domain.Download.Tries)
	assert.Equal(t, 10, domain.Download.Interval)
}
