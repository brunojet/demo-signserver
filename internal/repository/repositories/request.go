package repositories

import (
	"context"
	"log"

	"demo-signserver/internal/repository/domain"

	"demo-signserver/pkg/config"
	db_services "demo-signserver/pkg/repository/services"
)

const (
	REQUEST_RESOURCE_NAME = "request_table"
)

type RequestRepository struct {
	Dynamo *db_services.DynamoDBService
}

func NewRequestRepositoryCustom(configInstance string) *RequestRepository {
	cfg := config.GetConfigInstance(configInstance)
	resource := cfg.GetResource(REQUEST_RESOURCE_NAME)
	dynamo, err := db_services.NewDynamoDBService(resource.Name, db_services.ID_KEY, "")
	if err != nil {
		log.Fatalf("Erro ao inicializar DynamoDBService: %v", err)
	}
	return &RequestRepository{Dynamo: dynamo}
}

func NewRequestRepository() *RequestRepository {
	return NewRequestRepositoryCustom(config.DefaultInstanceName)
}

func (s *RequestRepository) CreateRequest(request *domain.SignRequest) error {
	return s.Dynamo.CreateItem(context.Background(), request)
}

func (s *RequestRepository) GetRequestByID(ID string) (*domain.SignRequest, error) {
	var request domain.SignRequest
	err := s.Dynamo.GetItem(context.Background(), ID, &request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (s *RequestRepository) UpdateRequest(ID string, request *domain.SignRequest) error {
	return s.Dynamo.UpdateItem(context.Background(), ID, request)
}
