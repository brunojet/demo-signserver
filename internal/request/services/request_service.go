package services

import (
	"demo-signserver/internal/repository/domain"
	"demo-signserver/internal/repository/repositories"
	"demo-signserver/internal/request/dtos"
	"errors"
)

type S3Service interface {
	GeneratePresignedURL(bucket, key string) (string, error)
}

type mockS3Service struct{}

func (m *mockS3Service) GeneratePresignedURL(bucket, key string) (string, error) {
	return "https://mock-s3-url/" + bucket + "/" + key, nil
}

type RequestService struct {
	service   *repositories.SignRequestService
	s3Service S3Service
}

var request_service *repositories.SignRequestService = nil
var s3_service S3Service = &mockS3Service{}

func SetRequestServiceMock(mock *repositories.SignRequestService) {
	request_service = mock
}

func SetS3ServiceMock(mock S3Service) {
	s3_service = mock
}

func getRequestService() *repositories.SignRequestService {
	if request_service == nil {
		request_service = repositories.NewSignRequestService()
	}
	return request_service
}

func getS3Service() S3Service {
	return s3_service
}

func NewRequestService() *RequestService {
	return &RequestService{service: getRequestService(), s3Service: getS3Service()}
}

func (s *RequestService) CreateRequest(request *domain.SignRequest) (string, error) {
	profileRepo := repositories.NewSignProfileService()
	profile, err := profileRepo.GetProfileByID(*request.SignerProfileId)
	if err != nil || profile == nil {
		return "", errors.New("profile_id não encontrado")
	}
	return s.service.CreateRequest(request)
}

func (s *RequestService) GetRequestByID(id string) (*domain.SignRequest, error) {
	return s.service.GetRequestByID(id)
}

func (s *RequestService) GetSignerStatusByID(id string) (*dtos.GetResponseDTO, error) {
	record, err := s.service.GetRequestByID(id)
	if (err != nil) || (record == nil) {
		return nil, err
	}
	response := dtos.NewGetResponseDTOFromDomain(record, getPresignedUrlFromDomain(record.SignedFile))
	return response, err
}

func getPresignedUrlFromDomain(bucketInfo *domain.BucketInfo) *string {
	if bucketInfo == nil || bucketInfo.BucketName == "" || bucketInfo.ObjectKey == "" {
		return nil
	}
	url, err := getS3Service().GeneratePresignedURL(bucketInfo.BucketName, bucketInfo.ObjectKey)
	if err != nil {
		return nil
	}
	return &url
}
