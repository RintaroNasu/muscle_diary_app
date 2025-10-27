package repository

import (
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"gorm.io/gorm"
)

type SummaryRepository interface {
	CountTrainingDays(userID uint) (int64, error)
	GetLatestWeight(userID uint) (*float64, *time.Time, error)
	GetProfileBasics(userID uint) (*float64, *float64, error)
}

type summaryRepository struct {
	db *gorm.DB
}

func NewSummaryRepository(db *gorm.DB) SummaryRepository {
	return &summaryRepository{db: db}
}

func (r *summaryRepository) CountTrainingDays(userID uint) (int64, error) {
	var cnt int64
	err := r.db.
		Model(&models.WorkoutRecord{}).
		Where("user_id = ?", userID).
		Distinct("trained_on").
		Count(&cnt).Error
	return cnt, err
}

func (r *summaryRepository) GetLatestWeight(userID uint) (*float64, *time.Time, error) {
	type row struct {
		BodyWeight float64
		TrainedOn  time.Time
	}

	var out row
	tx := r.db.
		Model(&models.WorkoutRecord{}).
		Where("user_id = ?", userID).
		Select("body_weight, trained_on").
		Order("trained_on DESC, id DESC").
		Limit(1).
		Scan(&out)
	if tx.Error != nil {
		return nil, nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, nil, nil
	}
	return &out.BodyWeight, &out.TrainedOn, nil
}

func (r *summaryRepository) GetProfileBasics(userID uint) (*float64, *float64, error) {
	var u models.User
	if err := r.db.
		Select("id, height, goal_weight").
		First(&u, userID).Error; err != nil {
		return nil, nil, err
	}
	return u.Height, u.GoalWeight, nil
}
