package domain

import "demo-signserver/pkg/repository/domain"

// ManufacturerType representa os fabricantes suportados para assinatura.
type Signer string

const (
	SignerPositivo Signer = "positivo"
	SignerGertec   Signer = "gertec"
)

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
	SignStepCreated           SignStep = "created"
	SignStepUploaded          SignStep = "uploaded"
	SignStepSigning           SignStep = "signing"
	SignStepSigned            SignStep = "signed"
	SignStepDownloadRequested SignStep = "download_requested"
	SignStepSigningFailed     SignStep = "signing_failed"
)

// IntentError representa um erro ocorrido em um passo do fluxo.
type IntentError struct {
	Code    string `dynamodbav:"code"`
	Message string `dynamodbav:"message"`
}

// IntentHistoryEntry representa um registro de histórico de um passo do fluxo.
type RequestHistoryEntry struct {
	Timestamp int64        `dynamodbav:"timestamp"`
	SignStep  *SignStep    `dynamodbav:"sign_step,omitempty"`
	Error     *IntentError `dynamodbav:"error"`
}

// SignRequest representa a entidade de intenção de assinatura.
type SignRequest struct {
	domain.BaseDomain
	SignProfileId string                `dynamodbav:"sign_profile_id"`
	SigningStatus SignStep              `dynamodbav:"signing_status"`
	UnsignedFile  BucketInfo            `dynamodbav:"unsigned_file"`
	SignedFile    BucketInfo            `dynamodbav:"signed_file"`
	WebhookURL    string                `dynamodbav:"webhook_url"`
	History       []RequestHistoryEntry `dynamodbav:"history"`
}

// TransferInfo representa informações de transferência (upload/download) de arquivos.
type TransferInfo struct {
	URL      string `dynamodbav:"url"`
	Tries    int    `dynamodbav:"tries"`
	Interval int    `dynamodbav:"interval"` // em segundos
}

// SignProfile representa o perfil de dispositivo associado à intenção.
type SignProfile struct {
	domain.BaseDomain
	Signer      *Signer                `dynamodbav:"signer,omitempty"`
	ProfileId   *string                `dynamodbav:"profile_id,omitempty"`
	Description *string                `dynamodbav:"description,omitempty"`
	Configs     *[]DeviceProfileConfig `dynamodbav:"configs,omitempty"`
	Upload      *TransferInfo          `dynamodbav:"upload,omitempty"`
	Download    *TransferInfo          `dynamodbav:"download,omitempty"`
}

// IntentRepository define o contrato para operações de persistência de intents.
type IntentRepository interface {
	CreateIntent(intent *SignRequest) error
	GetIntentByID(id string) (*SignRequest, error)
	UpdateIntent(intent *SignRequest) error
}

// DeviceProfileRepository define o contrato para operações de persistência de perfis de dispositivo.
type DeviceProfileRepository interface {
	CreateDeviceProfile(profile *SignProfile) error
	GetDeviceProfileByID(id string) (*SignProfile, error)
	UpdateDeviceProfile(profile *SignProfile) error
}

// NOTA DE ARQUITETURA:
//
// As interfaces IntentRepository e DeviceProfileRepository definem contratos genéricos para persistência de entidades de domínio.
// As implementações concretas dessas interfaces (adapters para DynamoDB, Postgres, etc.) devem residir em pkg/repository/services/ (código de produção) e pkg/repository/mock/ (mocks para testes).
//
// Os mocks de infraestrutura (ex: MockDynamoDBClient) devem ser mantidos em pkg/repository/mock, garantindo que apenas código de teste dependa deles.
//
// O mesmo padrão se aplica a storage: código de produção em pkg/storage/services e mocks em pkg/storage/mock.
//
// Um service/fachada em pkg/repository/services pode ser responsável por inicializar o adapter correto (conforme configuração) e expor métodos genéricos como CreateIntent, UpdateIntent, etc.,
// recebendo DTOs ou entidades de domínio e delegando para o adapter implementado.
//
// O preenchimento de PK, SK, timestamps e regras de negócio deve ser feito na camada de aplicação (ex: internal/intent/application/intent_service.go),
// mantendo o domínio limpo e a infraestrutura desacoplada.
//
// Essa abordagem permite centralizar a persistência, facilitar testes, troca de backend e reuso entre diferentes contextos (intent, signing, etc.).
//
// Exemplo de estrutura:
//   internal/
//     repository/
//       domain/      # interfaces de domínio
//   pkg/
//     repository/
//       services/    # implementações concretas (DynamoDB, Postgres, etc.)
//       mock/        # mocks para testes
//     storage/
//       services/    # serviços de storage (S3, STS, etc.)
//       mock/        # mocks de storage para testes
//
// Dessa forma, garantimos separação clara entre domínio, infraestrutura e testes, evitando dependências acidentais de mocks em produção e facilitando manutenção e evolução do projeto.
