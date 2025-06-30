package domain

import (
	"demo-signserver/pkg/repository/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignerError(t *testing.T) {
	err := &SignerError{Code: "E1", Message: "msg"}
	assert.Equal(t, "E1", err.Code)
	assert.Equal(t, "msg", err.Message)
}

func TestRequestHistoryEntry(t *testing.T) {
	step := SignerStepCreated
	hist := RequestHistoryEntry{
		Timestamp:  123456,
		SignerStep: &step,
		Error:      &SignerError{Code: "E2", Message: "fail"},
	}
	assert.Equal(t, int64(123456), hist.Timestamp)
	assert.Equal(t, SignerStepCreated, *hist.SignerStep)
	assert.Equal(t, "E2", hist.Error.Code)
}

func TestSignRequest_Fields(t *testing.T) {
	profileID := "pid"
	status := SignerStepSigned
	file := &BucketInfo{BucketName: "b", ObjectKey: "o", Size: 1, SHA256: "s"}
	hist := []RequestHistoryEntry{{Timestamp: 1}}
	// signer := "" // Remove this line since Signer expects *Signer, not *string
	var signer *Signer = nil // or initialize a *Signer if needed
	sr := &SignRequest{
		BaseDomain:      domain.BaseDomain{},
		Signer:          signer,
		SignerProfileId: &profileID,
		SignerStatus:    &status,
		UnsignedFile:    file,
		SignedFile:      file,
		WebhookURL:      nil,
		History:         &hist,
	}
	assert.NotNil(t, sr.BaseDomain)
	assert.Nil(t, sr.Signer)
	assert.Equal(t, &profileID, sr.SignerProfileId)
	assert.Equal(t, &status, sr.SignerStatus)
	assert.Equal(t, file, sr.UnsignedFile)
	assert.Equal(t, file, sr.SignedFile)
	assert.Equal(t, &hist, sr.History)
	assert.Nil(t, sr.WebhookURL) // BaseDomain should set an ID
}
