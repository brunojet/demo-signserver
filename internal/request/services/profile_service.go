package services

import (
	"demo-signserver/internal/config"
	"demo-signserver/internal/repository/domain"
	"demo-signserver/internal/repository/repositories"
)

type ProfileService struct {
	repository *repositories.ProfileRepository
}

func NewProfileService() *ProfileService {
	cfg := config.GetSignServerConfig()
	repository := repositories.NewProfileRepository(cfg.ProfileTableName)
	return &ProfileService{repository: repository}
}

func (s *ProfileService) CreateProfile(profile *domain.SignerProfile) error {
	return s.repository.CreateProfile(profile)
}

func (s *ProfileService) GetProfileByID(id string) (*domain.SignerProfile, error) {
	return s.repository.GetProfileByID(id)
}

func (s *ProfileService) UpdateProfile(id string, profile *domain.SignerProfile) error {
	return s.repository.UpdateProfile(id, profile)
}
