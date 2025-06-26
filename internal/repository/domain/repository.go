package domain

// ManufacturerType representa os fabricantes suportados para assinatura.
type ManufacturerType string

const (
	ManufacturerPositivo ManufacturerType = "positivo"
	ManufacturerGertec   ManufacturerType = "gertec"
)

// BaseEntity contém campos comuns para entidades persistidas.
type BaseEntity struct {
	PK        string `dynamodbav:"pk"`
	SK        string `dynamodbav:"sk"`
	CreatedAt int64  `dynamodbav:"created_at"`
	UpdatedAt int64  `dynamodbav:"updated_at"`
}

// DeviceProfileConfig representa configurações específicas de perfil de dispositivo.
type DeviceProfileConfig struct {
	Key   string `dynamodbav:"key"`
	Value string `dynamodbav:"value"`
}

// BucketInfo representa informações de um arquivo em um bucket.
type BucketInfo struct {
	BucketName string `dynamodbav:"bucket_name"`
	ObjectKey  string `dynamodbav:"object_key"`
	Size       int64  `dynamodbav:"size"`
	SHA256     string `dynamodbav:"sha256"`
}

// SignStep representa os possíveis passos do fluxo de assinatura.
type SignStep string

const (
	SignStepCreated       SignStep = "created"
	SignStepUploaded      SignStep = "uploaded"
	SignStepSigning       SignStep = "signing"
	SignStepSigned        SignStep = "signed"
	SignStepSigningFailed SignStep = "signing_failed"
)

// NotificationStep representa os possíveis passos do fluxo de notificação.
type NotificationStep string

const (
	NotificationStepNotifying          NotificationStep = "notifying"
	NotificationStepNotified           NotificationStep = "notified"
	NotificationStepNotificationFailed NotificationStep = "notification_failed"
)

// IntentError representa um erro ocorrido em um passo do fluxo.
type IntentError struct {
	Code    string `dynamodbav:"code"`
	Message string `dynamodbav:"message"`
}

// IntentHistoryEntry representa um registro de histórico de um passo do fluxo.
type IntentHistoryEntry struct {
	Timestamp        int64             `dynamodbav:"timestamp"`
	SignStep         *SignStep         `dynamodbav:"sign_step,omitempty"`
	NotificationStep *NotificationStep `dynamodbav:"notification_step,omitempty"`
	Error            *IntentError      `dynamodbav:"error"`
}

// Intent representa a entidade de intenção de assinatura.
type Intent struct {
	BaseEntity
	SigningStatus      SignStep             `dynamodbav:"signing_status"`
	NotificationStatus NotificationStep     `dynamodbav:"notification_status"`
	UnsignedFile       BucketInfo           `dynamodbav:"unsigned_file"`
	SignedFile         BucketInfo           `dynamodbav:"signed_file"`
	WebhookURL         string               `dynamodbav:"webhook_url"`
	History            []IntentHistoryEntry `dynamodbav:"history"`
}

// TransferInfo representa informações de transferência (upload/download) de arquivos.
type TransferInfo struct {
	URL      string `dynamodbav:"url"`
	Tries    int    `dynamodbav:"tries"`
	Interval int    `dynamodbav:"interval"` // em segundos
}

// DeviceProfile representa o perfil de dispositivo associado à intenção.
type DeviceProfile struct {
	BaseEntity
	Name         string                `dynamodbav:"name"`
	Manufacturer ManufacturerType      `dynamodbav:"manufacturer"`
	Configs      []DeviceProfileConfig `dynamodbav:"configs"`
	Upload       TransferInfo          `dynamodbav:"upload"`
	Download     TransferInfo          `dynamodbav:"download"`
}

// IntentRepository define o contrato para operações de persistência de intents.
type IntentRepository interface {
	CreateIntent(intent *Intent) error
	GetIntentByID(id string) (*Intent, error)
	UpdateIntent(intent *Intent) error
}

// DeviceProfileRepository define o contrato para operações de persistência de perfis de dispositivo.
type DeviceProfileRepository interface {
	CreateDeviceProfile(profile *DeviceProfile) error
	GetDeviceProfileByID(id string) (*DeviceProfile, error)
	UpdateDeviceProfile(profile *DeviceProfile) error
}
