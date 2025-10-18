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
	FindByIDAndUserID(id uint, userID uint) (*models.WorkoutRecord, error)
	Update(record *models.WorkoutRecord) error
	Delete(id uint, userID uint) error
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

func (r *workoutRepository) FindByIDAndUserID(id uint, userID uint) (*models.WorkoutRecord, error) {
	var record models.WorkoutRecord
	err := r.db.
		Preload("Sets").
		Preload("Exercise").
		Where("id = ? AND user_id = ?", id, userID).
		First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *workoutRepository) Update(record *models.WorkoutRecord) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().
			Where("workout_record_id = ?", record.ID).
			Delete(&models.WorkoutSet{}).Error; err != nil {
			return err
		}

		if err := tx.
			Omit("Sets.*").
			Model(&models.WorkoutRecord{}).
			Where("id = ?", record.ID).
			Updates(map[string]any{
				"body_weight": record.BodyWeight,
				"exercise_id": record.ExerciseID,
				"trained_on":  record.TrainedOn,
			}).Error; err != nil {
			return err
		}

		for i := range record.Sets {
			record.Sets[i].ID = 0
			record.Sets[i].WorkoutRecordID = record.ID
		}
		if len(record.Sets) > 0 {
			if err := tx.Create(&record.Sets).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *workoutRepository) Delete(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Unscoped().Delete(&models.WorkoutRecord{}).Error
}
