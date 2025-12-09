package repository

import (
	"time"

	"gorm.io/gorm"
)

type TimelineItem struct {
	RecordID     uint
	UserID       uint
	UserEmail    string
	ExerciseName string
	BodyWeight   float64
	TrainedOn    time.Time
	Comment      string
}

type TimelineRepository interface {
	FindPublicRecords() ([]TimelineItem, error)
}

type timelineRepository struct {
	db *gorm.DB
}

func NewTimelineRepository(db *gorm.DB) TimelineRepository {
	return &timelineRepository{db: db}
}

func (r *timelineRepository) FindPublicRecords() ([]TimelineItem, error) {
	var rows []TimelineItem

	err := r.db.
		Table("workout_records").
		Select(`
			workout_records.id          AS record_id,
			workout_records.user_id     AS user_id,
			users.email                 AS user_email,
			exercises.name              AS exercise_name,
			workout_records.body_weight AS body_weight,
			workout_records.trained_on  AS trained_on,
			workout_records.comment     AS comment
		`).
		Joins("JOIN users ON users.id = workout_records.user_id").
		Joins("JOIN exercises ON exercises.id = workout_records.exercise_id").
		Where("workout_records.is_public = ?", true).
		Order("workout_records.trained_on DESC, workout_records.id DESC").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	return rows, nil
}
