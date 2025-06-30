package application

import (
	"demo-signserver/internal/repository/domain"
	"demo-signserver/internal/repository/services"
)

type ProfileService struct {
	service *services.SignProfileService
}

var profile_service *services.SignProfileService = nil

func SetProfileServiceMock(mock *services.SignProfileService) {
	profile_service = mock
}

func getProfileService() *services.SignProfileService {
	if profile_service == nil {
		profile_service = services.NewSignProfileService()
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
