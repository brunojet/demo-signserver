package services

import (
	"demo-signserver/internal/config"
	"demo-signserver/internal/repository/domain"
	"demo-signserver/internal/repository/repositories"
	"demo-signserver/pkg/storage/adapters"
	"errors"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type RequestService struct {
	repository *repositories.RequestRepository
	storage    adapters.StorageServiceInterface
}

func NewRequestService() *RequestService {
	cfg := config.GetSignServerConfig()
	methods := config.GetSignServerMethods()
	storage := methods.NewStorageService(cfg.StorageBucketName)
	repository := repositories.NewRequestRepository()
	return &RequestService{repository: repository, storage: storage}
}

func (s *RequestService) CreateRequest(request *domain.SignRequest) (*domain.SignRequestResponse, error) {
	cfg := config.GetSignServerConfig()
	profileRepo := repositories.NewProfileRepository(cfg.ProfileTableName)
	profile, err := profileRepo.GetProfileByID(*request.SignerProfileId)
	if err != nil || profile == nil {
		return nil, errors.New("profile_id não encontrado")
	}

	ID, url, err := s.generatePresignedPutURL()

	if err != nil || ID == "" {
		return nil, errors.New("erro ao gerar URL pré-assinada")
	}

	request.SetID(ID)
	err = s.repository.CreateRequest(request)

	if err != nil {
		return nil, err
	}

	response := &domain.SignRequestResponse{
		ID:           ID,
		SignerStatus: *request.SignerStatus,
		SignerError:  request.GetLastError(),
		HttpMethod:   domain.HttpMethodPut,
		UploadURL:    url,
	}

	return response, err
}

func (s *RequestService) GetRequestByID(id string) (*domain.SignRequest, error) {
	return s.repository.GetRequestByID(id)
}

func (s *RequestService) GetSignerStatusByID(id string) (*domain.SignGetResponse, error) {
	record, err := s.repository.GetRequestByID(id)
	if (err != nil) || (record == nil) {
		return nil, err
	}
	response := &domain.SignGetResponse{
		ID:           record.ID,
		SignerStatus: *record.SignerStatus,
		SignerError:  record.GetLastError(),
	}

	s.getPresignedGetUrl(response, record.SignedFile)

	return response, err
}

func (s *RequestService) getPresignedGetUrl(response *domain.SignGetResponse, bucketInfo *domain.BucketInfo) {
	if response == nil || bucketInfo == nil || bucketInfo.Key == "" {
		return
	}
	url, err := s.storage.GeneratePresignedURL(adapters.HttpMethodGet, bucketInfo.Key, 15*time.Minute)
	if err != nil {
		return
	}
	httpMethod := domain.HttpMethodGet
	response.HttpMethod = &httpMethod
	response.DownloadURL = &url
}

// Gera um nome de arquivo único, gera URL pré-assinada e retorna (nome, url, erro)
func (s *RequestService) generatePresignedPutURL() (string, string, error) {
	fileName := uuid.New().String()
	url, err := s.storage.GeneratePresignedURL(adapters.HttpMethodPut, filepath.Join("unsigned", fileName), 15*time.Minute)
	if err != nil {
		return "", "", err
	}
	return fileName, url, nil
}
