package models

import (
	"time"

	"gorm.io/gorm"
)

type WorkoutRecord struct {
	gorm.Model
	UserID     uint
	ExerciseID uint `gorm:"foreignKey:ExerciseID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	BodyWeight float64
	TrainedAt  time.Time
	Sets       []WorkoutSet `gorm:"constraint:OnDelete:CASCADE"`
}

type WorkoutSet struct {
	gorm.Model
	WorkoutRecordID uint
	SetNo           int
	Reps            int
	ExerciseWeight  float64
}
