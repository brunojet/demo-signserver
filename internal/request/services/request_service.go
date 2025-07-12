package services

import (
	"demo-signserver/internal/config"
	"demo-signserver/internal/repository/domain"
	"demo-signserver/internal/repository/repositories"
	"demo-signserver/pkg/storage"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type RequestService struct {
	repository  *repositories.RequestRepository
	storage     storage.StorageAdapter
	storagePath string
}

func NewRequestService() *RequestService {
	cfg := config.GetSignServerConfig()
	methods := config.GetSignServerMethods()
	storage := methods.NewStorageService()
	repository := repositories.NewRequestRepository()
	return &RequestService{repository: repository, storage: storage, storagePath: cfg.StorageBucketName}
}

func (s *RequestService) CreateRequest(request *domain.SignRequest) (*domain.SignRequestResponse, error) {

	profileRepo := repositories.NewProfileRepository()
	_, err := profileRepo.GetProfileByID(*request.SignerProfileId)

	if err != nil {
		return nil, err
	}

	request.SetID(uuid.New().String())

	request.UnsignedFile = &storage.FileInfo{
		StoragePath: s.storagePath,
		FilePath:    filepath.Join("unsigned", request.ID),
	}

	request.SignedFile = &storage.FileInfo{
		StoragePath: s.storagePath,
		FilePath:    filepath.Join("signed", request.ID),
	}

	url, err := s.generatePresignedPutURL(request.UnsignedFile)

	if err != nil {
		return nil, err
	}

	err = s.repository.CreateRequest(request)

	if err != nil {
		return nil, err
	}

	return &domain.SignRequestResponse{
		ID:           request.ID,
		SignerStatus: *request.SignerStatus,
		HttpMethod:   domain.HttpMethodPut,
		UploadURL:    url,
	}, nil
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

func (s *RequestService) getPresignedGetUrl(response *domain.SignGetResponse, fileInfo *storage.FileInfo) {
	if response == nil {
		return
	}
	url, err := s.storage.GeneratePresignedURL(storage.HttpMethodGet, fileInfo, 15*time.Minute)
	if err != nil {
		return
	}
	httpMethod := domain.HttpMethodGet
	response.HttpMethod = httpMethod
	response.DownloadURL = url
}

func (s *RequestService) generatePresignedPutURL(fileInfo *storage.FileInfo) (string, error) {
	url, err := s.storage.GeneratePresignedURL(storage.HttpMethodPut, fileInfo, 15*time.Minute)
	if err != nil {
		return "", err
	}
	return url, nil
}
