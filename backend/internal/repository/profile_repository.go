package repository

import (
	"errors"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"gorm.io/gorm"
)

type ProfileRepository interface {
	GetProfile(userID uint) (*models.User, error)
	UpdateProfile(userID uint, height *float64, goalWeight *float64) error
}

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) GetProfile(userID uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *profileRepository) UpdateProfile(userID uint, height *float64, goalWeight *float64) error {
	updates := map[string]interface{}{
		"height":      height,
		"goal_weight": goalWeight,
	}
	result := r.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
