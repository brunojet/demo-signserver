package services

import (
	"demo-signserver/internal/repository/domain"
	"demo-signserver/internal/repository/repositories"
)

type ProfileService struct {
	service *repositories.SignProfileService
}

var profile_service *repositories.SignProfileService = nil

func SetProfileServiceMock(mock *repositories.SignProfileService) {
	profile_service = mock
}

func getProfileService() *repositories.SignProfileService {
	if profile_service == nil {
		profile_service = repositories.NewSignProfileService()
	}
	return profile_service
}

func NewProfileService() *ProfileService {
	return &ProfileService{service: getProfileService()}
}

func (s *ProfileService) CreateProfile(profile *domain.SignerProfile) (string, error) {
	return s.service.CreateProfile(profile)
}

func (s *ProfileService) GetProfileByID(id string) (*domain.SignerProfile, error) {
	return s.service.GetProfileByID(id)
}

func (s *ProfileService) UpdateProfile(id string, profile *domain.SignerProfile) error {
	return s.service.UpdateProfile(id, profile)
}
