package services

import (
	"demo-signserver/internal/repository/domain"
	"demo-signserver/internal/repository/repositories"
	"demo-signserver/internal/request/dtos"
	storage_services "demo-signserver/pkg/storage/services"
	"errors"
	"os"
	"time"
)

type RequestService struct {
	service   *repositories.SignRequestService
	s3Service *storage_services.S3Service
}

var request_service *repositories.SignRequestService = nil
var s3_service *storage_services.S3Service = nil

func SetRequestServiceMock(mock *repositories.SignRequestService) {
	request_service = mock
}

func SetS3ServiceMock(mock *storage_services.S3Service) {
	s3_service = mock
}

func getRequestService() *repositories.SignRequestService {
	if request_service == nil {
		request_service = repositories.NewSignRequestService()
	}
	return request_service
}

func getS3Service() *storage_services.S3Service {
	if s3_service == nil {
		bucket := os.Getenv("SIGN_STORAGE_BUCKET")
		realS3, err := storage_services.NewS3Service(bucket)
		if err != nil {
			panic("Erro ao criar S3Service real: " + err.Error())
		}
		s3_service = realS3
	}
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
	url, err := getS3Service().GeneratePresignedURL(bucketInfo.ObjectKey, 15*time.Minute)
	if err != nil {
		return nil
	}
	return &url
}
