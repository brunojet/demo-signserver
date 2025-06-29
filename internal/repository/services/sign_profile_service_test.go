package services

import (
	"demo-signserver/internal/repository/domain"
	"fmt"
	"os"
	"testing"
)

func TestSignProfileService_CRUD(t *testing.T) {
	os.Setenv("PROJECT_NAME", "signserver")
	os.Setenv("ENVIRONMENT", "dev")
	os.Setenv("SIGN_PROFILE_TABLE", "profile")

	service := NewSignProfileService()

	profile := &domain.SignProfile{
		Description: "Dispositivos Postivo perfil 001",
		Signer:      domain.SignerPositivo,
		ProfileId:   "002",
		Configs:     []domain.DeviceProfileConfig{{Key: "k", Value: "v"}},
		Upload:      domain.TransferInfo{URL: "https://example.com/upload", Tries: 3, Interval: 5},
		Download:    domain.TransferInfo{URL: "https://example.com/download", Tries: 3, Interval: 5},
	}

	err := service.CreateProfile(profile)
	if err != nil {
		t.Fatalf("Erro ao criar perfil: %v", err)
	}

	fetched, err := service.GetProfileByID(fmt.Sprintf("%s#%s", profile.Signer, profile.ProfileId))
	if err != nil {
		t.Fatalf("Erro ao buscar perfil: %v", err)
	}
	if fetched.Description != profile.Description {
		t.Errorf("Nome esperado %s, obtido %s", profile.Description, fetched.Description)
	}

	fetched.Description = "Unit Test Profile Updated"
	err = service.UpdateProfile(fetched)
	if err != nil {
		t.Fatalf("Erro ao atualizar perfil: %v", err)
	}

	fetched, err = service.GetProfileByID(fmt.Sprintf("%s#%s", fetched.Signer, fetched.ProfileId))
	if err != nil {
		t.Fatalf("Erro ao buscar perfil atualizado: %v", err)
	}
	if fetched.Description != "Unit Test Profile Updated" {
		t.Errorf("Nome esperado 'Unit Test Profile Updated', obtido %s", fetched.Description)
	}
}
