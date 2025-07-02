package dtos

import (
	"demo-signserver/internal/repository/domain"
)

type CreateSignRequestDTO struct {
	ProfileId  string  `json:"profile_id" binding:"required"`
	WebhookURL *string `json:"callback_url,omitempty" binding:"omitempty,url"`
}

type CreateSignResponseDTO struct {
	ID        string `json:"id" binding:"required"`
	UploadURL string `json:"pre_signed_url" binding:"required,url"`
}

func (d *CreateSignRequestDTO) GetDomainCreateSignRequest() *domain.SignRequest {
	domainReq := domain.SignRequest{
		SignerProfileId: &d.ProfileId,
		WebhookURL:      d.WebhookURL,
	}

	domainReq.SetSignerStatus(domain.SignerStepCreated, nil)
	return &domainReq

}

type GetResponseDTO struct {
	ID           string              `json:"id" binding:"required"`
	SignerStatus domain.SignerStep   `json:"signer_status,omitempty" binding:"required"`
	SignerError  *domain.SignerError `json:"signer_error,omitempty"`
	DownloadURL  *string             `json:"pre_signed_url" binding:"required,url"`
}

func NewGetResponseDTOFromDomain(req *domain.SignRequest, downloadUrl *string) *GetResponseDTO {
	if req == nil {
		return nil
	}

	return &GetResponseDTO{
		ID:           req.ID,
		SignerStatus: *req.SignerStatus,
		SignerError:  req.GetLastError(),
		DownloadURL:  downloadUrl,
	}
}
