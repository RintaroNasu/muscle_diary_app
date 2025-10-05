package repository

import (
	"context"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"gorm.io/gorm"
)

type ExerciseRepository interface {
	List(ctx context.Context) ([]models.Exercise, error)
}

type exerciseRepository struct {
	db *gorm.DB
}

func NewExerciseRepository(db *gorm.DB) ExerciseRepository {
	return &exerciseRepository{db: db}
}

func (r *exerciseRepository) List(ctx context.Context) ([]models.Exercise, error) {
	var xs []models.Exercise
	if err := r.db.WithContext(ctx).
		Select("id", "name").
		Order("id ASC").
		Find(&xs).Error; err != nil {
		return nil, err
	}
	return xs, nil
}
