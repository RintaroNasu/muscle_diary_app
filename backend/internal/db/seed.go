package db

import (
	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"gorm.io/gorm"
)

var defalutExercises = []models.Exercise{
	{Name: "ベンチプレス"},
	{Name: "スクワット"},
	{Name: "デッドリフト"},
}

func Seed(db *gorm.DB) error {
	for _, exercise := range defalutExercises {
		var count int64

		if err := db.Model(&models.Exercise{}).Where("name = ?", exercise.Name).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := db.Create(&exercise).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
