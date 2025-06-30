package dtos

import "demo-signserver/internal/repository/domain"

type TransferInfoDTO struct {
	Url      string `json:"url" binding:"required,url"`
	Tries    int    `json:"tries" binding:"required,min=1,max=10"`
	Interval int    `json:"interval" binding:"required,min=5,max=30"`
}

type CreateSignerProfileDTO struct {
	Signer      domain.Signer           `json:"signer" binding:"required,oneof=positivo gertec"`
	ProfileId   string                  `json:"profile_id" binding:"required,len=3,numeric"`
	Description string                  `json:"description" binding:"required,min=3,max=255"`
	Configs     *[]domain.ProfileConfig `json:"configs,omitempty"`
	Upload      TransferInfoDTO         `json:"upload" binding:"required"`
	Download    TransferInfoDTO         `json:"download" binding:"required"`
}

type UpdateSignerProfileDTO struct {
	Description *string                 `json:"description,omitempty" binding:"min=3,max=255"`
	Configs     *[]domain.ProfileConfig `json:"configs,omitempty"`
	Upload      *TransferInfoDTO        `json:"upload,omitempty"`
	Download    *TransferInfoDTO        `json:"download,omitempty"`
}

type GetSignerProfileDTO = domain.SignerProfile
