package service

import (
	"fmt"
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
)

type WorkoutService interface {
	CreateWorkoutRecord(userID uint, bodyWeight float64, exerciseID uint, trainedAt time.Time, sets []WorkoutSetData) (*models.WorkoutRecord, error)
}

type WorkoutSetData struct {
	SetNo          int
	Reps           int
	ExerciseWeight float64
}

func NewWorkoutService(repo repository.WorkoutRepository) WorkoutService {
	return &workoutService{repo: repo}
}

type workoutService struct {
	repo repository.WorkoutRepository
}

func (s *workoutService) CreateWorkoutRecord(userID uint, bodyWeight float64, exerciseID uint, trainedAt time.Time, sets []WorkoutSetData) (*models.WorkoutRecord, error) {
	record := &models.WorkoutRecord{
		UserID:       userID,
		ExerciseID:   exerciseID,
		BodyWeight:   bodyWeight,
		TrainedAt:    trainedAt,
	}

	for _, setData := range sets {
		set := models.WorkoutSet{
			SetNo:          setData.SetNo,
			Reps:           setData.Reps,
			ExerciseWeight: setData.ExerciseWeight,
		}
		record.Sets = append(record.Sets, set)
	}

	if err := s.repo.Create(record); err != nil {
		return nil, fmt.Errorf("create workout record failed: %w", err)
	}

	return record, nil
}
