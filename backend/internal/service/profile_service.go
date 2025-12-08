package service

import (
	"errors"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
)

type ProfileService interface {
	GetProfile(userID uint) (*models.User, error)
	UpdateProfile(userID uint, height *float64, goalWeight *float64) (*models.User, error)
}

type profileService struct {
	repo repository.ProfileRepository
}

func NewProfileService(repo repository.ProfileRepository) ProfileService {
	return &profileService{repo: repo}
}

func (s *profileService) GetProfile(userID uint) (*models.User, error) {
	user, err := s.repo.GetProfile(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *profileService) UpdateProfile(userID uint, height *float64, goalWeight *float64) (*models.User, error) {
	if err := s.repo.UpdateProfile(userID, height, goalWeight); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return s.GetProfile(userID)
}
