package repository

import (
	"fmt"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(u *models.User) error
	FindByEmail(email string) (*models.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) Create(u *models.User) error {
	if err := r.db.Create(u).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, fmt.Errorf("failed to find user by email=%s: %w", email, err)
	}
	return &u, nil
}
