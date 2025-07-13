package dtos

import "demo-signserver/internal/repository/domain"

type TransferInfoDTO struct {
	Url      string `json:"url" binding:"required,url"`
	Tries    int    `json:"tries" binding:"required,min=1,max=10"`
	Interval int    `json:"interval" binding:"required,min=5,max=30"`
}

type CreateProfileDTO struct {
	Signer      domain.Signer           `json:"signer" binding:"required,oneof=positivo gertec"`
	ProfileId   string                  `json:"profile_id" binding:"required,len=3,numeric"`
	Description string                  `json:"description" binding:"required,min=3,max=255"`
	Configs     *[]domain.ProfileConfig `json:"configs,omitempty"`
	ContentType string                  `json:"content_type" binding:"required,oneof=binary/octet-stream application/vnd.android.package-archive"`
	Upload      TransferInfoDTO         `json:"upload" binding:"required"`
	Download    TransferInfoDTO         `json:"download" binding:"required"`
}

type UpdateProfileDTO struct {
	Description *string                 `json:"description,omitempty" binding:"omitempty,min=3,max=255"`
	Configs     *[]domain.ProfileConfig `json:"configs,omitempty"`
	ContentType *string                 `json:"content_type,omitempty" binding:"omitempty,oneof=binary/octet-stream application/vnd.android.package-archive"`
	Upload      *TransferInfoDTO        `json:"upload,omitempty"`
	Download    *TransferInfoDTO        `json:"download,omitempty"`
}

type GetSignerProfileDTO = domain.SignerProfile

func (d *CreateProfileDTO) GetDomainCreateSignerProfile() *domain.SignerProfile {
	return &domain.SignerProfile{
		Signer:      &d.Signer,
		ProfileId:   &d.ProfileId,
		Description: &d.Description,
		Configs:     d.Configs,
		Upload:      convertToDomainTransferInfo(&d.Upload),
		Download:    convertToDomainTransferInfo(&d.Download),
	}
}

func (d *UpdateProfileDTO) GetDomainUpdateSignerProfile() *domain.SignerProfile {
	return &domain.SignerProfile{
		Description: d.Description,
		Configs:     d.Configs,
		Upload:      convertToDomainTransferInfo(d.Upload),
		Download:    convertToDomainTransferInfo(d.Download),
	}
}

func convertToDomainTransferInfo(dto *TransferInfoDTO) *domain.TransferInfo {
	if dto == nil {
		return nil
	}
	return &domain.TransferInfo{
		URL:      dto.Url,
		Tries:    dto.Tries,
		Interval: dto.Interval,
	}
}
