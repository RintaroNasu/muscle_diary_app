package service

import (
	"errors"

	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
)

type WorkoutLikeService interface {
	Like(userID uint, recordID uint) error
	Unlike(userID uint, recordID uint) error
}

type workoutLikeService struct {
	repo repository.WorkoutLikeRepository
}

func NewWorkoutLikeService(repo repository.WorkoutLikeRepository) WorkoutLikeService {
	return &workoutLikeService{repo: repo}
}

func (s *workoutLikeService) Like(userID uint, recordID uint) error {
	public, err := s.repo.IsRecordPublic(recordID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrRecordNotFound
		}
		return err
	}
	if !public {
		return ErrForbiddenPrivateRecord
	}

	return s.repo.CreateLike(userID, recordID)
}

func (s *workoutLikeService) Unlike(userID uint, recordID uint) error {
	public, err := s.repo.IsRecordPublic(recordID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrRecordNotFound
		}
		return err
	}
	if !public {
		return ErrForbiddenPrivateRecord
	}

	return s.repo.DeleteLike(userID, recordID)
}
