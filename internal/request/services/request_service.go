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
	s3_service := getS3Service(bucket)
	return &RequestService{service: getRequestService(), s3Service: s3_service}
}

func (s *RequestService) CreateRequest(request *domain.SignRequest) (*domain.SignRequestResponse, error) {
	profileRepo := repositories.NewSignProfileService()
	profile, err := profileRepo.GetProfileByID(*request.SignerProfileId)
	if err != nil || profile == nil {
		return nil, errors.New("profile_id não encontrado")
	}

	ID, url, err := s.generatePresignedPutURL()

	if err != nil || ID == "" {
		return nil, errors.New("erro ao gerar URL pré-assinada")
	}

	request.SetUnsingedBucketInfo(s.s3Service.Bucket, ID)
	request.SetID(ID)
	err = s.service.CreateRequest(request)

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
	}

	s.getPresignedGetUrl(response, record.SignedFile)

	return response, err
}

func (s *RequestService) getPresignedGetUrl(response *domain.SignGetResponse, bucketInfo *domain.BucketInfo) {
	if response == nil || bucketInfo == nil || bucketInfo.ObjectKey == "" {
		return
	}
	url, err := s.s3Service.GeneratePresignedGetURL(bucketInfo.ObjectKey, 15*time.Minute)
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
	url, err := s.s3Service.GeneratePresignedPutURL(fileName, 15*time.Minute)
	if err != nil {
		return "", "", err
	}
	return fileName, url, nil
}
