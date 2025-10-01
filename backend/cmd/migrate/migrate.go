package migrate

import (
	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"gorm.io/gorm"
)

func Migrate(conn *gorm.DB) error {
	return conn.AutoMigrate(
		&models.User{},
		&models.WorkoutRecord{},
		&models.WorkoutSet{},
	)
}
