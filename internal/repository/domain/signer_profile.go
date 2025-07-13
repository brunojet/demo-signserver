package domain

import "demo-signserver/pkg/repository/domain"

// ManufacturerType representa os fabricantes suportados para assinatura.
type Signer string

const (
	SignerPositivo Signer = "positivo"
	SignerGertec   Signer = "gertec"
)

// ProfileConfig representa configurações específicas de perfil de dispositivo.
type ProfileConfig struct {
	Key   string `dynamodbav:"key"`
	Value string `dynamodbav:"value"`
}

// TransferInfo representa informações de transferência (upload/download) de arquivos.
type TransferInfo struct {
	URL      string `dynamodbav:"url"`
	Tries    int    `dynamodbav:"tries"`
	Interval int    `dynamodbav:"interval"` // em segundos
}

// SignerProfile representa o perfil de dispositivo associado à intenção.
type SignerProfile struct {
	domain.BaseDomain
	Signer      *Signer          `dynamodbav:"signer,omitempty"`
	ProfileId   *string          `dynamodbav:"profile_id,omitempty"`
	Description *string          `dynamodbav:"description,omitempty"`
	Configs     *[]ProfileConfig `dynamodbav:"configs,omitempty"`
	Upload      *TransferInfo    `dynamodbav:"upload,omitempty"`
	Download    *TransferInfo    `dynamodbav:"download,omitempty"`
}
