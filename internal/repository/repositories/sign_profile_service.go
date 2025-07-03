package repositories

import (
	"context"
	"demo-signserver/internal/repository/domain"
	"log"
	"os"

	db_services "demo-signserver/pkg/repository/services"
)

const SIGNER_KEY = "signer"
const PROFILE_ID_KEY = "profile_id"

type SignProfileServiceInterface interface {
	CreateProfile(profile *domain.SignerProfile) (string, error)
	GetProfileByID(ID string) (*domain.SignerProfile, error)
	UpdateProfile(ID string, profile *domain.SignerProfile) error
}
type SignProfileService struct {
	Dynamo *db_services.DynamoDBService
}

func NewSignProfileService() *SignProfileService {
	table := os.Getenv("SIGN_PROFILE_TABLE")
	dynamo, err := db_services.NewDynamoDBService(table, SIGNER_KEY, PROFILE_ID_KEY)
	if err != nil {
		log.Fatalf("Erro ao inicializar DynamoDBService: %v", err)
	}
	return &SignProfileService{Dynamo: dynamo}
}

func (s *SignProfileService) CreateProfile(profile *domain.SignerProfile) error {
	return s.Dynamo.CreateItem(context.Background(), profile)
}

func (s *SignProfileService) GetProfileByID(ID string) (*domain.SignerProfile, error) {
	var profile domain.SignerProfile
	err := s.Dynamo.GetItem(context.Background(), ID, &profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (s *SignProfileService) UpdateProfile(ID string, profile *domain.SignerProfile) error {
	return s.Dynamo.UpdateItem(context.Background(), ID, profile)
}
