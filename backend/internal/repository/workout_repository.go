package repository

import (
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"gorm.io/gorm"
)

type WorkoutRepository interface {
	Create(record *models.WorkoutRecord) error
	FindByUserAndDay(userID uint, day time.Time) ([]models.WorkoutRecord, error)
	FindRecordDaysInMonth(userID uint, year int, month int) ([]time.Time, error)
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

func (r *workoutRepository) FindByUserAndDay(userID uint, day time.Time) ([]models.WorkoutRecord, error) {
	var records []models.WorkoutRecord
	err := r.db.
		Where("user_id = ? AND trained_on = ?", userID, day).Preload("Exercise").Preload("Sets").
		Order("id ASC").
		Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r *workoutRepository) FindRecordDaysInMonth(userID uint, year int, month int) ([]time.Time, error) {
	var dates []time.Time
	err := r.db.
		Model(&models.WorkoutRecord{}).
		Where("user_id = ? AND EXTRACT(YEAR FROM trained_on) = ? AND EXTRACT(MONTH FROM trained_on) = ?",
			userID, year, month).
		Distinct("trained_on").
		Order("trained_on ASC").
		Pluck("trained_on", &dates).Error

	if err != nil {
		return nil, err
	}
	return dates, nil
}
