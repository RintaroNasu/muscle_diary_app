package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
	"gorm.io/gorm"
)

type WorkoutService interface {
	CreateWorkoutRecord(userID uint, bodyWeight float64, exerciseID uint, trainedOn time.Time, sets []WorkoutSetData) (*models.WorkoutRecord, error)
	GetDailyRecords(userID uint, day time.Time) ([]models.WorkoutRecord, error)
	GetMonthRecordDays(userID uint, year int, month int) ([]time.Time, error)
}

type WorkoutSetData struct {
	SetNo          int
	Reps           int
	ExerciseWeight float64
}

type workoutService struct {
	repo repository.WorkoutRepository
}

var (
	ErrNoSets           = errors.New("no sets")
	ErrInvalidSetValue  = errors.New("invalid set value")
	ErrExerciseNotFound = errors.New("exercise not found")
)

func NewWorkoutService(repo repository.WorkoutRepository) WorkoutService {
	return &workoutService{repo: repo}
}

func (s *workoutService) CreateWorkoutRecord(userID uint, bodyWeight float64, exerciseID uint, trainedOn time.Time, sets []WorkoutSetData) (*models.WorkoutRecord, error) {
	if len(sets) == 0 {
		return nil, ErrNoSets
	}

	for _, st := range sets {
		if st.SetNo <= 0 || st.Reps <= 0 || st.ExerciseWeight < 0 {
			return nil, ErrInvalidSetValue
		}
	}

	record := &models.WorkoutRecord{
		UserID:     userID,
		ExerciseID: exerciseID,
		BodyWeight: bodyWeight,
		TrainedOn:  trainedOn,
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
		if errors.Is(err, gorm.ErrForeignKeyViolated) || strings.Contains(err.Error(), "foreign key") {
			return nil, ErrExerciseNotFound
	}
	return nil, fmt.Errorf("create workout record failed: %w", err)
	}

	return record, nil
}

func (s *workoutService) GetDailyRecords(userID uint, day time.Time) ([]models.WorkoutRecord, error) {
	records, err := s.repo.FindByUserAndDay(userID, day)
	if err != nil {
		return nil, fmt.Errorf("fetch daily records failed: %w", err)
	}
	return records, nil
}

func (s *workoutService) GetMonthRecordDays(userID uint, year int, month int) ([]time.Time, error) {
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("invalid month: %d", month)
	}

	days, err := s.repo.FindRecordDaysInMonth(userID, year, month)
	if err != nil {
		return nil, fmt.Errorf("fetch month record days failed: %w", err)
	}
	return days, nil
}
