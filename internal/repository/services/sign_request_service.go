package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"demo-signserver/internal/repository/domain"

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

type SignRequestService struct {
	Dynamo *db_services.DynamoDBService
}

func NewSignRequestService() *SignRequestService {
	project := os.Getenv("PROJECT_NAME")
	env := os.Getenv("ENVIRONMENT")
	table_name := os.Getenv("SIGN_REQUEST_TABLE")
	table := fmt.Sprintf("%s-%s-%s", project, env, table_name)
	dynamo, err := db_services.NewDynamoDBService(table, "", "")
	if err != nil {
		log.Fatalf("Erro ao inicializar DynamoDBService: %v", err)
	}
	return &SignRequestService{Dynamo: dynamo}
}

func (s *SignRequestService) CreateIntent(intent *domain.SignRequest) error {
	item, err := db_services.MarshalItem(intent)
	if err != nil {
		return err
	}
	return s.Dynamo.CreateItem(context.TODO(), item)
}

func (s *SignRequestService) GetIntentByID(ID string) (*domain.SignRequest, error) {
	var intent domain.SignRequest
	err := s.Dynamo.GetItem(context.TODO(), ID, &intent)
	if err != nil {
		return nil, err
	}
	return &intent, nil
}

func (s *SignRequestService) UpdateIntent(intent *domain.SignRequest) error {
	item, err := db_services.MarshalItem(intent)
	if err != nil {
		return err
	}
	return s.Dynamo.UpdateItem(context.TODO(), item)
}
