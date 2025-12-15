package repository

import (
	"errors"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"gorm.io/gorm"
)

type WorkoutLikeRepository interface {
	CreateLike(userID uint, recordID uint) error
	DeleteLike(userID uint, recordID uint) error
	IsRecordPublic(recordID uint) (bool, error)
	IsLikedByMe(userID uint, recordID uint) (bool, error)
}

type workoutLikeRepository struct {
	db *gorm.DB
}

func NewWorkoutLikeRepository(db *gorm.DB) WorkoutLikeRepository {
	return &workoutLikeRepository{db: db}
}

func (r *workoutLikeRepository) CreateLike(userID uint, recordID uint) error {
	like := models.WorkoutLike{
		UserID:   userID,
		RecordID: recordID,
	}

	if err := r.db.Create(&like).Error; err != nil {
		liked, e := r.IsLikedByMe(userID, recordID)
		if e == nil && liked {
			return nil
		}
		return err
	}

	return nil
}

func (r *workoutLikeRepository) DeleteLike(userID uint, recordID uint) error {
	res := r.db.
		Unscoped().
		Where("user_id = ? AND record_id = ?", userID, recordID).
		Delete(&models.WorkoutLike{})

	if res.Error != nil {
		return res.Error
	}

	// rows=0 でも冪等で成功扱い
	return nil
}

func (r *workoutLikeRepository) IsRecordPublic(recordID uint) (bool, error) {
	var rec models.WorkoutRecord
	if err := r.db.Select("id, is_public").First(&rec, recordID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, ErrNotFound
		}
		return false, err
	}
	return rec.IsPublic, nil
}

func (r *workoutLikeRepository) IsLikedByMe(userID uint, recordID uint) (bool, error) {
	var like models.WorkoutLike
	err := r.db.
		Where("user_id = ? AND record_id = ?", userID, recordID).
		First(&like).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
