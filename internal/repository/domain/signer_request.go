package domain

import "demo-signserver/pkg/repository/domain"

// BucketInfo representa informações de um arquivo em um bucket.
type BucketInfo struct {
	BucketName string `dynamodbav:"bucket_name"`
	ObjectKey  string `dynamodbav:"object_key"`
	Size       int64  `dynamodbav:"size"`
	SHA256     string `dynamodbav:"sha256"`
}

// SignerStep representa os possíveis passos do fluxo de assinatura.
type SignerStep string

const (
	SignerStepCreated           SignerStep = "created"
	SignerStepUploaded          SignerStep = "uploaded"
	SignerStepSigning           SignerStep = "signing"
	SignerStepSigned            SignerStep = "signed"
	SignerStepDownloadRequested SignerStep = "download_requested"
	SignerStepSigningFailed     SignerStep = "signing_failed"
)

// SignerError representa um erro ocorrido em um passo do fluxo.
type SignerError struct {
	Code    string `dynamodbav:"code"`
	Message string `dynamodbav:"message"`
}

// IntentHistoryEntry representa um registro de histórico de um passo do fluxo.
type RequestHistoryEntry struct {
	Timestamp  int64        `dynamodbav:"timestamp"`
	SignerStep *SignerStep  `dynamodbav:"sign_step,omitempty"`
	Error      *SignerError `dynamodbav:"error"`
}

// SignRequest representa a entidade de intenção de assinatura.
type SignRequest struct {
	domain.BaseDomain
	Signer          *Signer                `dynamodbav:"signer,omitempty"`
	SignerProfileId *string                `dynamodbav:"signer_profile_id,omitempty"`
	SignerStatus    *SignerStep            `dynamodbav:"signer_status,omitempty"`
	UnsignedFile    *BucketInfo            `dynamodbav:"unsigned_file,omitempty"`
	SignedFile      *BucketInfo            `dynamodbav:"signed_file,omitempty"`
	WebhookURL      *string                `dynamodbav:"webhook_url,omitempty"`
	History         *[]RequestHistoryEntry `dynamodbav:"history,omitempty"`
}
