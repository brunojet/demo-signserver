package repositories

import (
	"context"
	"demo-signserver/internal/config"
	"demo-signserver/internal/repository/domain"

	db_services "demo-signserver/pkg/repository/services"
)

const (
	REQUEST_RESOURCE_NAME = "request_table"
)

type RequestRepository struct {
	Dynamo *db_services.DynamoDBService
}

func NewRequestRepository() *RequestRepository {
	cfg := config.GetSignServerConfig()
	methods := config.GetSignServerMethods()
	dynamo := methods.NewDynamoDBService(cfg.RequestTableName, db_services.ID_KEY, db_services.NO_KEY)
	return &RequestRepository{Dynamo: dynamo}
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
