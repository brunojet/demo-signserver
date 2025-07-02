package repositories

import (
	"context"
	"demo-signserver/internal/repository/domain"
	db_services "demo-signserver/pkg/repository/services"
	"log"
	"os"
	"testing"
)

func init() {
	os.Setenv("PROJECT_NAME", "signserver")
	os.Setenv("ENVIRONMENT", "dev")
	os.Setenv("SIGN_PROFILE_TABLE", "profile")
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8001")
	os.Setenv("AWS_ACCESS_KEY_ID", "fake")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "fake")
	db := db_services.NewDB("Dummy")
	if err := db.DeleteTable(context.Background(), "profile"); err != nil {
		log.Fatalf("Error deleting table request: %v", err)
	}
	if err := db.CreateTable(context.Background(), "profile", "profile_id"); err != nil {
		log.Fatalf("Error creating table request: %v", err)
	}
}

func TestSignProfileService_CRUD(t *testing.T) {
	service := NewSignProfileService()

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

	ID, err := service.CreateProfile(profile)
	if err != nil {
		t.Fatalf("Erro ao criar perfil: %v", err)
	}

	fetched, err := service.GetProfileByID(ID)
	if err != nil {
		t.Fatalf("Erro ao buscar perfil: %v", err)
	}
	if *fetched.Description != *profile.Description {
		t.Errorf("Nome esperado %s, obtido %s", *profile.Description, *fetched.Description)
	}

	updatedDesc := "Unit Test Profile Updated"
	update := &domain.SignerProfile{
		Description: &updatedDesc,
	}

	err = service.UpdateProfile(ID, update)
	if err != nil {
		t.Fatalf("Erro ao atualizar perfil: %v", err)
	}

	fetched, err = service.GetProfileByID(ID)
	if err != nil {
		t.Fatalf("Erro ao buscar perfil atualizado: %v", err)
	}
	if fetched.Description == nil || update.Description == nil || *fetched.Description != *update.Description {
		t.Errorf("Nome esperado '%s', obtido %v", *update.Description, fetched.Description)
	}
}
