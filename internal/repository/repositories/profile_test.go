package repositories

import (
	"demo-signserver/internal/config"
	"demo-signserver/internal/repository/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignProfileService_CRUD(t *testing.T) {
	cfg := config.GetSignServerConfig()
	service := NewProfileRepository(cfg.ProfileTableName)

	desc := "Dispositivos Postivo perfil 001"
	profileId := "015"
	signer := domain.SignerPositivo
	profile := &domain.SignerProfile{
		Description: &desc,
		Signer:      &signer,
		ProfileId:   &profileId,
		Configs:     &[]domain.ProfileConfig{{Key: "k", Value: "v"}},
		Upload:      &domain.TransferInfo{URL: "https://example.com/upload", Tries: 3, Interval: 5},
		Download:    &domain.TransferInfo{URL: "https://example.com/download", Tries: 3, Interval: 5},
	}

	err := service.CreateProfile(profile)
	assert.NoError(t, err)

	fetched, err := service.GetProfileByID(profile.ID)
	assert.NoError(t, err)
	if *fetched.Description != *profile.Description {
		t.Errorf("Nome esperado %s, obtido %s", *profile.Description, *fetched.Description)
	}

	updatedDesc := "Unit Test Profile Updated"
	update := &domain.SignerProfile{
		Description: &updatedDesc,
	}

	err = service.UpdateProfile(profile.ID, update)
	assert.NoError(t, err)
	fetched, err = service.GetProfileByID(profile.ID)
	if err != nil {
		t.Fatalf("Erro ao buscar perfil atualizado: %v", err)
	}
	if fetched.Description == nil || update.Description == nil || *fetched.Description != *update.Description {
		t.Errorf("Nome esperado '%s', obtido %v", *update.Description, fetched.Description)
	}
}
