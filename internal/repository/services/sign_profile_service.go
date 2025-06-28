package services

import (
	"context"
	"demo-signserver/internal/repository/domain"
	"fmt"
	"log"
	"os"

	db_services "demo-signserver/pkg/repository/services"

	"github.com/joho/godotenv"
)

func init() {
	// Carrega o arquivo .env da pasta config
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}
}

type SignProfileService struct {
	Dynamo *db_services.DynamoDBService
}

func NewSignProfileService() *SignProfileService {
	project := os.Getenv("PROJECT_NAME")
	env := os.Getenv("ENVIRONMENT")
	table_name := os.Getenv("SIGN_PROFILE_TABLE")
	table := fmt.Sprintf("%s-%s-%s", project, env, table_name)
	dynamo, err := db_services.NewDynamoDBService(table)
	if err != nil {
		log.Fatalf("Erro ao inicializar DynamoDBService: %v", err)
	}
	return &SignProfileService{Dynamo: dynamo}
}

func (s *SignProfileService) CreateProfile(profile *domain.SignProfile) error {
	db_services.SetPKSKFromID(profile)
	db_services.SetTimestamps(profile, true)
	item, err := db_services.MarshalItem(profile)
	if err != nil {
		return err
	}
	return s.Dynamo.PutItem(context.TODO(), item)
}

func (s *SignProfileService) GetProfileByID(id string) (*domain.SignProfile, error) {
	key := db_services.BuildKeyFromID(id)
	item, err := s.Dynamo.GetItem(context.TODO(), key)
	if err != nil {
		return nil, err
	}
	var profile domain.SignProfile
	err = db_services.UnmarshalItem(item, &profile)
	if err != nil {
		return nil, err
	}
	db_services.SetIDFromPKSK(&profile)
	return &profile, nil
}

func (s *SignProfileService) UpdateProfile(profile *domain.SignProfile) error {
	db_services.SetPKSKFromID(profile)
	db_services.SetTimestamps(profile, false)
	item, err := db_services.MarshalItem(profile)
	if err != nil {
		return err
	}
	return s.Dynamo.PutItem(context.TODO(), item)
}
