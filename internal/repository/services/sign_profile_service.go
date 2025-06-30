package services

import (
	"context"
	"demo-signserver/internal/repository/domain"
	"log"
	"os"

	db_services "demo-signserver/pkg/repository/services"
)

const SIGNER_KEY = "signer"
const PROFILE_ID_KEY = "profile_id"

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

func (s *SignProfileService) CreateProfile(profile *domain.SignerProfile) (string, error) {
	return s.Dynamo.CreateItem(context.TODO(), profile)
}

func (s *SignProfileService) GetProfileByID(ID string) (*domain.SignerProfile, error) {
	var profile domain.SignerProfile
	err := s.Dynamo.GetItem(context.TODO(), ID, &profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (s *SignProfileService) UpdateProfile(ID string, profile *domain.SignerProfile) error {
	return s.Dynamo.UpdateItem(context.TODO(), ID, profile)
}
