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
	TrainedOn  time.Time    `gorm:"type:date;not null;index"`
	Sets       []WorkoutSet `gorm:"constraint:OnDelete:CASCADE"`
	Exercise   Exercise     `gorm:"foreignKey:ExerciseID"`
	IsPublic   bool         `gorm:"default:false"`
	Comment    string       `gorm:"type:text"`
}

type WorkoutSet struct {
	gorm.Model
	WorkoutRecordID uint
	SetNo           int
	Reps            int
	ExerciseWeight  float64
}
