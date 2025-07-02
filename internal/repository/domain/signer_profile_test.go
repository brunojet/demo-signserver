package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfileConfig(t *testing.T) {
	cfg := ProfileConfig{Key: "k", Value: "v"}
	assert.Equal(t, "k", cfg.Key)
	assert.Equal(t, "v", cfg.Value)
}

func TestTransferInfo(t *testing.T) {
	tr := TransferInfo{URL: "u", Tries: 2, Interval: 10}
	assert.Equal(t, "u", tr.URL)
	assert.Equal(t, 2, tr.Tries)
	assert.Equal(t, 10, tr.Interval)
}

func TestSignerProfile_Fields(t *testing.T) {
	profileID := "pid"
	desc := "desc"
	signer := SignerPositivo
	cfgs := []ProfileConfig{{Key: "a", Value: "b"}}
	upload := &TransferInfo{URL: "up", Tries: 1, Interval: 5}
	download := &TransferInfo{URL: "down", Tries: 2, Interval: 10}
	sp := &SignerProfile{
		Signer:      &signer,
		ProfileId:   &profileID,
		Description: &desc,
		Configs:     &cfgs,
		Upload:      upload,
		Download:    download,
	}
	assert.Equal(t, &signer, sp.Signer)
	assert.Equal(t, &profileID, sp.ProfileId)
	assert.Equal(t, &desc, sp.Description)
	assert.Equal(t, &cfgs, sp.Configs)
	assert.Equal(t, upload, sp.Upload)
	assert.Equal(t, download, sp.Download)
}
