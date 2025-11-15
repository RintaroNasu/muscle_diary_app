package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"gorm.io/gorm"
)

type RankingRepository interface {
	MonthlyGymDays(ctx context.Context, from, to time.Time) ([]GymDaysRow, error)
}

type rankingRepository struct {
	db *gorm.DB
}

func NewRankingRepository(db *gorm.DB) RankingRepository {
	return &rankingRepository{db: db}
}

type GymDaysRow struct {
	UserID            uint
	Email             string
	TotalTrainingDays int64
}

type TotalVolumeRow struct {
	UserID      uint
	Email       string
	TotalVolume float64
}

func (r *rankingRepository) MonthlyGymDays(ctx context.Context, from, to time.Time) ([]GymDaysRow, error) {
	var rows []GymDaysRow
	err := r.db.WithContext(ctx).
		Model(&models.WorkoutRecord{}).
		Select("workout_records.user_id AS user_id, users.email AS email,COUNT(DISTINCT workout_records.trained_on) AS total_training_days").
		Joins("JOIN users ON workout_records.user_id = users.id").
		Where("trained_on >= ? AND trained_on < ?", from, to).
		Group("workout_records.user_id, users.email").
		Order("total_training_days DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	fmt.Println("rows", rows)
	return rows, nil
}
