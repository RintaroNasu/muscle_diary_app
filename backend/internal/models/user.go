package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email      string `gorm:"unique"`
	Password   string
	Records    []WorkoutRecord `gorm:"constraint:OnDelete:CASCADE"`
	Height     *float64        `gorm:"type:numeric(4,1)"`
	GoalWeight *float64        `gorm:"type:numeric(4,1)"`
}
