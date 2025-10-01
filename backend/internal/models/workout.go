package models

import (
	"time"

	"gorm.io/gorm"
)

type WorkoutRecord struct {
	gorm.Model
	UserID       uint
	ExerciseName string
	BodyWeight   float64
	TrainedAt    time.Time
	Sets         []WorkoutSet `gorm:"constraint:OnDelete:CASCADE"`
}

type WorkoutSet struct {
	gorm.Model
	WorkoutRecordID uint
	SetNo           int
	Reps            int
	ExerciseWeight  float64
}
