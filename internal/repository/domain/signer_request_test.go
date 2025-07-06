package domain

import (
	"demo-signserver/pkg/repository/domain"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignerError(t *testing.T) {
	err := &SignerError{Code: "E1", Message: "msg"}
	assert.Equal(t, "E1", err.Code)
	assert.Equal(t, "msg", err.Message)
}

func TestRequestHistoryEntry(t *testing.T) {
	step := SignerStatusCreated
	hist := RequestHistoryEntry{
		CreatedAt:    "123456",
		SignerStatus: &step,
		Error:        &SignerError{Code: "E2", Message: "fail"},
	}
	assert.Equal(t, "123456", hist.CreatedAt)
	assert.Equal(t, SignerStatusCreated, *hist.SignerStatus)
	assert.Equal(t, "E2", hist.Error.Code)
}

func TestSignRequest_Fields(t *testing.T) {
	profileID := "pid"
	status := SignerStatusSigned
	file := &BucketInfo{BucketName: "b", ObjectKey: "o", Size: 1, SHA256: "s"}
	hist := []RequestHistoryEntry{{CreatedAt: "1"}}
	// signer := "" // Remove this line since Signer expects *Signer, not *string
	sr := &SignRequest{
		BaseDomain:      domain.BaseDomain{},
		SignerProfileId: &profileID,
		SignerStatus:    &status,
		UnsignedFile:    file,
		SignedFile:      file,
		WebhookURL:      nil,
		History:         &hist,
	}
	assert.NotNil(t, sr.BaseDomain)
	assert.Equal(t, &profileID, sr.SignerProfileId)
	assert.Equal(t, &status, sr.SignerStatus)
	assert.Equal(t, file, sr.UnsignedFile)
	assert.Equal(t, file, sr.SignedFile)
	assert.Equal(t, &hist, sr.History)
	assert.Nil(t, sr.WebhookURL) // BaseDomain should set an ID
}

func TestSignRequest_SetSignerStatus_NewHistoryAndStatus(t *testing.T) {
	sr := &SignRequest{}
	step := SignerStatusSigning
	err := SignerError{Code: "E3", Message: "error msg"}

	sr.SetSignerStatus(step, &err)

	assert.NotNil(t, sr.SignerStatus)
	assert.Equal(t, step, *sr.SignerStatus)
	assert.NotNil(t, sr.History)
	assert.Equal(t, 1, len(*sr.History))
	last := (*sr.History)[0]
	assert.Equal(t, step, *last.SignerStatus)
	assert.Equal(t, err.Code, last.Error.Code)
	assert.Equal(t, err.Message, last.Error.Message)
}

func TestSignRequest_SetSignerStatus_AppendHistory(t *testing.T) {
	sr := &SignRequest{}
	firstStep := SignerStatusCreated
	firstErr := SignerError{Code: "E1", Message: "msg1"}
	sr.SetSignerStatus(firstStep, &firstErr)

	secondStep := SignerStatusSigned
	secondErr := SignerError{Code: "E2", Message: "msg2"}
	sr.SetSignerStatus(secondStep, &secondErr)

	assert.NotNil(t, sr.SignerStatus)
	assert.Equal(t, secondStep, *sr.SignerStatus)
	assert.NotNil(t, sr.History)
	assert.Equal(t, 2, len(*sr.History))
	last := (*sr.History)[1]
	assert.Equal(t, secondStep, *last.SignerStatus)
	assert.Equal(t, secondErr.Code, last.Error.Code)
	assert.Equal(t, secondErr.Message, last.Error.Message)
}

func TestSignRequest_GetLastError(t *testing.T) {
	errMsg := "erro de assinatura"
	errCode := "E123"
	err := &SignerError{Code: errCode, Message: errMsg}
	stepOk := SignerStatusSigned
	stepFail := SignerStatusSigningFailed

	history := []RequestHistoryEntry{
		{SignerStatus: &stepOk, Error: nil},
		{SignerStatus: &stepFail, Error: err},
	}

	req := &SignRequest{
		History: &history,
	}

	got := req.GetLastError()
	if !reflect.DeepEqual(got, err) {
		t.Errorf("Expected error %+v, got %+v", err, got)
	}

	// Testa sem erro
	req2 := &SignRequest{History: &[]RequestHistoryEntry{{SignerStatus: &stepOk, Error: nil}}}
	if req2.GetLastError() != nil {
		t.Errorf("Expected nil, got %+v", req2.GetLastError())
	}

	// Testa com history vazio
	req3 := &SignRequest{History: &[]RequestHistoryEntry{}}
	if req3.GetLastError() != nil {
		t.Errorf("Expected nil, got %+v", req3.GetLastError())
	}

	// Testa com history nil
	req4 := &SignRequest{}
	if req4.GetLastError() != nil {
		t.Errorf("Expected nil, got %+v", req4.GetLastError())
	}
}
