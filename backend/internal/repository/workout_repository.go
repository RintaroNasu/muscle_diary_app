package repository

import (
	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"gorm.io/gorm"
)

type WorkoutRepository interface {
	Create(record *models.WorkoutRecord) error
}

type workoutRepository struct {
	db *gorm.DB
}

func NewWorkoutRepository(db *gorm.DB) WorkoutRepository {
	return &workoutRepository{db: db}
}

func (r *workoutRepository) Create(record *models.WorkoutRecord) error {
	return r.db.Create(record).Error
}
