package domain

import (
	"demo-signserver/pkg/repository/domain"
	"time"
)

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

type HttpMethod string

const (
	HttpMethodPut HttpMethod = "PUT"
	HttpMethodGet HttpMethod = "GET"
)

// SignerError representa um erro ocorrido em um passo do fluxo.
type SignerError struct {
	Code    string `json:"code" dynamodbav:"code"`
	Message string `json:"message" dynamodbav:"message"`
}

// IntentHistoryEntry representa um registro de histórico de um passo do fluxo.
type RequestHistoryEntry struct {
	CreatedAt  string       `json:"created_at" dynamodbav:"timestamp"`
	SignerStep *SignerStep  `json:"sign_step" dynamodbav:"sign_step,omitempty"`
	Error      *SignerError `json:"error" dynamodbav:"error"`
}

type SignerRequestInerface interface {
	domain.BaseDomainInterface
	SetSignerStatus(step SignerStep, err SignerError)
}

// SignRequest representa a entidade de intenção de assinatura.
type SignRequest struct {
	domain.BaseDomain
	SignerProfileId *string                `json:"profile_id,omitempty" dynamodbav:"signer_profile_id,omitempty"`
	SignerStatus    *SignerStep            `json:"signer_status,omitempty" dynamodbav:"signer_status,omitempty"`
	UnsignedFile    *BucketInfo            `json:"unsigned_file,omitempty" dynamodbav:"unsigned_file,omitempty"`
	SignedFile      *BucketInfo            `json:"signed_file,omitempty" dynamodbav:"signed_file,omitempty"`
	WebhookURL      *string                `json:"webhook_url,omitempty" dynamodbav:"webhook_url,omitempty"`
	History         *[]RequestHistoryEntry `json:"history,omitempty" dynamodbav:"history,omitempty"`
}

type SignRequestResponse struct {
	ID           string       `json:"id" binding:"required"`
	SignerStatus SignerStep   `json:"signer_status,omitempty" binding:"required"`
	SignerError  *SignerError `json:"signer_error,omitempty"`
	HttpMethod   HttpMethod   `json:"method" binding:"required,oneof=PUT"`
	UploadURL    string       `json:"upload_url" binding:"required,url"`
}

type SignGetResponse struct {
	ID           string       `json:"id" binding:"required"`
	SignerStatus SignerStep   `json:"signer_status,omitempty" binding:"required"`
	SignerError  *SignerError `json:"signer_error,omitempty"`
	HttpMethod   HttpMethod   `json:"method" binding:"required,oneof=GET"`
	DownloadURL  string       `json:"download_url" binding:"required,url"`
}

func (s *SignRequest) SetUnsingedBucketInfo(bucketName, objectKey string) {
	s.UnsignedFile = &BucketInfo{
		BucketName: bucketName,
		ObjectKey:  objectKey,
		Size:       0,
		SHA256:     "",
	}
}

func (s *SignRequest) SetSignerStatus(step SignerStep, err *SignerError) {
	if s.SignerStatus == nil {
		s.SignerStatus = &step
	} else {
		*s.SignerStatus = step
	}

	if s.History == nil {
		s.History = &[]RequestHistoryEntry{}
	}

	history := RequestHistoryEntry{
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
		SignerStep: &step,
		Error:      err,
	}
	*s.History = append(*s.History, history)
}

func (s *SignRequest) GetLastError() *SignerError {
	if s.History == nil || len(*s.History) == 0 {
		return nil
	}
	// Percorre do mais recente para o mais antigo
	for i := len(*s.History) - 1; i >= 0; i-- {
		h := (*s.History)[i]
		if h.SignerStep != nil && *h.SignerStep == SignerStepSigningFailed {
			return h.Error
		}
	}
	return nil
}
