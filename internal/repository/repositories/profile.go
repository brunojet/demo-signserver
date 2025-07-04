package repositories

import (
	"context"
	"demo-signserver/internal/repository/domain"
	"log"

	"demo-signserver/pkg/config"
	db_services "demo-signserver/pkg/repository/services"
)

const (
	PROFILE_RESOURCE_NAME = "profile_table"
	SIGNER_KEY            = "signer"
	PROFILE_ID_KEY        = "profile_id"
)

type ProfileRepositoryInterface interface {
	CreateProfile(profile *domain.SignerProfile) (string, error)
	GetProfileByID(ID string) (*domain.SignerProfile, error)
	UpdateProfile(ID string, profile *domain.SignerProfile) error
}
type ProfileRepository struct {
	Dynamo *db_services.DynamoDBService
}

func NewProfileRepositoryCustom(configInstance string) *ProfileRepository {
	cfg := config.GetConfigInstance(configInstance)
	resource := cfg.GetResource(PROFILE_RESOURCE_NAME)
	dynamo, err := db_services.NewDynamoDBService(resource.Name, SIGNER_KEY, PROFILE_ID_KEY)
	if err != nil {
		log.Fatalf("Erro ao inicializar DynamoDBService: %v", err)
	}
	return &ProfileRepository{Dynamo: dynamo}
}

func NewProfileRepository() *ProfileRepository {
	return NewProfileRepositoryCustom(config.DefaultInstanceName)
}

func (s *ProfileRepository) CreateProfile(profile *domain.SignerProfile) error {
	return s.Dynamo.CreateItem(context.Background(), profile)
}

func (s *ProfileRepository) GetProfileByID(ID string) (*domain.SignerProfile, error) {
	var profile domain.SignerProfile
	err := s.Dynamo.GetItem(context.Background(), ID, &profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (s *ProfileRepository) UpdateProfile(ID string, profile *domain.SignerProfile) error {
	return s.Dynamo.UpdateItem(context.Background(), ID, profile)
}
