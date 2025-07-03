package services

import (
	"demo-signserver/internal/repository/domain"
	"demo-signserver/internal/repository/repositories"
	storage_services "demo-signserver/pkg/storage/services"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
)

type RequestService struct {
	service   *repositories.SignRequestService
	s3Service *storage_services.S3Service
	bucket    string
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

func getS3Service(bucket string) *storage_services.S3Service {
	if s3_service == nil {
		realS3, err := storage_services.NewS3Service(bucket)
		if err != nil {
			panic("Erro ao criar S3Service real: " + err.Error())
		}
		s3_service = realS3
	}
	return s3_service
}

func NewRequestService() *RequestService {
	bucket := os.Getenv("SIGN_STORAGE_BUCKET")
	return &RequestService{service: getRequestService(), s3Service: getS3Service(bucket), bucket: bucket}
}

func (s *RequestService) CreateRequest(request *domain.SignRequest) (*domain.SignRequestResponse, error) {
	profileRepo := repositories.NewSignProfileService()
	profile, err := profileRepo.GetProfileByID(*request.SignerProfileId)
	if err != nil || profile == nil {
		return nil, errors.New("profile_id não encontrado")
	}

	ID, url, err := s.GenerateUniquePresignedURL()

	if err != nil || ID == "" {
		return nil, errors.New("erro ao gerar URL pré-assinada")
	}

	request.SetUnsingedBucketInfo(s.bucket, ID)
	request.SetID(ID)
	err = s.service.CreateRequest(request)

	if err != nil {
		return nil, err
	}

	response := &domain.SignRequestResponse{
		ID:           ID,
		SignerStatus: *request.SignerStatus,
		SignerError:  request.GetLastError(),
		UploadURL:    url,
	}

	return response, err
}

func (s *RequestService) GetRequestByID(id string) (*domain.SignRequest, error) {
	return s.service.GetRequestByID(id)
}

func (s *RequestService) GetSignerStatusByID(id string) (*domain.SignGetResponse, error) {
	record, err := s.service.GetRequestByID(id)
	if (err != nil) || (record == nil) {
		return nil, err
	}
	response := &domain.SignGetResponse{
		ID:           record.ID,
		SignerStatus: *record.SignerStatus,
		SignerError:  record.GetLastError(),
		DownloadURL:  s.getPresignedUrl(record.SignedFile),
	}
	return response, err
}

func (s *RequestService) getPresignedUrl(bucketInfo *domain.BucketInfo) string {
	if bucketInfo == nil || bucketInfo.BucketName == "" || bucketInfo.ObjectKey == "" {
		return ""
	}
	url, err := s.s3Service.GeneratePresignedURL(bucketInfo.ObjectKey, 15*time.Minute)
	if err != nil {
		return ""
	}
	return url
}

// Gera um nome de arquivo único, gera URL pré-assinada e retorna (nome, url, erro)
func (s *RequestService) GenerateUniquePresignedURL() (string, string, error) {
	fileName := uuid.New().String()
	url, err := s.s3Service.GeneratePresignedURL(fileName, 15*time.Minute)
	if err != nil {
		return "", "", err
	}
	return fileName, url, nil
}
