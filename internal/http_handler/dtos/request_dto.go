package dtos

import (
	"demo-signserver/internal/repository/domain"
)

type CreateSignRequestDTO struct {
	ProfileId  string  `json:"profile_id" binding:"required"`
	WebhookURL *string `json:"callback_url,omitempty" binding:"omitempty,url"`
}

type CreateSignResponseDTO struct {
	ID           string `json:"id" binding:"required"`
	PreSignedURL string `json:"pre_signed_url" binding:"required"`
}

func (d *CreateSignRequestDTO) GetDomainCreateSignRequest() *domain.SignRequest {
	signerStatus := domain.SignerStepCreated
	return &domain.SignRequest{
		SignerProfileId: &d.ProfileId,
		SignerStatus:    &signerStatus,
		WebhookURL:      d.WebhookURL,
	}
}

type GetSignRequestDTO = domain.SignRequest
