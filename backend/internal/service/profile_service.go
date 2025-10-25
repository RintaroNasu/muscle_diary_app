package service

import (
	"errors"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
	"gorm.io/gorm"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *profileService) UpdateProfile(userID uint, height *float64, goalWeight *float64) (*models.User, error) {
	if err := s.repo.UpdateProfile(userID, height, goalWeight); err != nil {
		return nil, err
	}

	// 更新後の最新データを返す
	return s.GetProfile(userID)
}
