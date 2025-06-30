package dtos

import "demo-signserver/internal/repository/domain"

type SignerProfileFieldsDTO struct {
	Description *string                 `json:"description,omitempty"`
	Configs     *[]domain.ProfileConfig `json:"configs,omitempty"`
	Upload      *domain.TransferInfo    `json:"upload,omitempty"`
	Download    *domain.TransferInfo    `json:"download,omitempty"`
}

type CreateSignerProfileDTO struct {
	Signer    *domain.Signer `json:"signer" binding:"required"`
	ProfileId *string        `json:"profile_id" binding:"required"`
	SignerProfileFieldsDTO
}

type UpdateSignerProfileDTO struct {
	SignerProfileFieldsDTO
}

type GetSignerProfileDTO = domain.SignerProfile
