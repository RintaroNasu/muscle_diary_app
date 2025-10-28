package service

import (
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
)

type HomeSummary struct {
	TotalTrainingDays int64      `json:"total_training_days"`
	LatestWeight      *float64   `json:"latest_weight,omitempty"`
	LatestTrainedOn   *time.Time `json:"trained_on,omitempty"`
	GoalWeight        *float64   `json:"goal_weight,omitempty"`
	Height            *float64   `json:"height,omitempty"`
}

type SummaryService interface {
	GetHomeSummary(userID uint) (*HomeSummary, error)
}

type summaryService struct {
	repo repository.SummaryRepository
}

func NewSummaryService(repo repository.SummaryRepository) SummaryService {
	return &summaryService{repo: repo}
}

func (s *summaryService) GetHomeSummary(userID uint) (*HomeSummary, error) {
	days, err := s.repo.CountTrainingDays(userID)
	if err != nil {
		return nil, err
	}

	latestW, latestOn, err := s.repo.GetLatestWeight(userID)
	if err != nil {
		return nil, err
	}

	height, goal, err := s.repo.GetProfileBasics(userID)
	if err != nil {
		return nil, err
	}

	return &HomeSummary{
		TotalTrainingDays: days,
		LatestWeight:      latestW,
		LatestTrainedOn:   latestOn,
		GoalWeight:        goal,
		Height:            height,
	}, nil
}
