package repositories

import (
	"context"
	"log"
	"os"

	"demo-signserver/internal/repository/domain"

	db_services "demo-signserver/pkg/repository/services"
)

type SignRequestService struct {
	Dynamo *db_services.DynamoDBService
}

func NewSignRequestService() *SignRequestService {
	table := os.Getenv("SIGN_REQUEST_TABLE")
	dynamo, err := db_services.NewDynamoDBService(table, "", "")
	if err != nil {
		log.Fatalf("Erro ao inicializar DynamoDBService: %v", err)
	}
	return &SignRequestService{Dynamo: dynamo}
}

func (s *SignRequestService) CreateRequest(request *domain.SignRequest) (string, error) {
	return s.Dynamo.CreateItem(context.Background(), request)
}

func (s *SignRequestService) GetRequestByID(ID string) (*domain.SignRequest, error) {
	var request domain.SignRequest
	err := s.Dynamo.GetItem(context.Background(), ID, &request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (s *SignRequestService) UpdateRequest(ID string, request *domain.SignRequest) error {
	return s.Dynamo.UpdateItem(context.Background(), ID, request)
}
