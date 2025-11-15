// internal/service/ranking_service.go
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
)

type RankingService interface {
	MonthlyGymDays(ctx context.Context, year, month int) ([]GymDaysDTO, error)
	MonthlyTotalVolume(ctx context.Context, year, month int) ([]TotalVolumeDTO, error)
}

type rankingService struct {
	repo repository.RankingRepository
}

func NewRankingService(repo repository.RankingRepository) RankingService {
	return &rankingService{repo: repo}
}

type GymDaysDTO struct {
	UserID            uint   `json:"user_id"`
	Email             string `json:"email"`
	TotalTrainingDays int64  `json:"total_training_days"`
}

type TotalVolumeDTO struct {
	UserID      uint    `json:"user_id"`
	Email       string  `json:"email"`
	TotalVolume float64 `json:"total_volume"`
}

func calcMonthRange(year, month int) (time.Time, time.Time, error) {
	if month < 1 || month > 12 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid month: %d", month)
	}
	if year < 2000 || year > 2100 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid year: %d", year)
	}

	loc, _ := time.LoadLocation("Asia/Tokyo")
	from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
	fmt.Println("from", from)
	to := from.AddDate(0, 1, 0)
	fmt.Println("to", to)
	return from, to, nil
}

func (s *rankingService) MonthlyGymDays(
	ctx context.Context, year, month int,
) ([]GymDaysDTO, error) {
	from, to, err := calcMonthRange(year, month)
	if err != nil {
		return nil, err
	}

	rows, err := s.repo.MonthlyGymDays(ctx, from, to)
	if err != nil {
		return nil, err
	}

	out := make([]GymDaysDTO, 0, len(rows))
	for _, r := range rows {
		out = append(out, GymDaysDTO{
			UserID:            r.UserID,
			Email:             r.Email,
			TotalTrainingDays: r.TotalTrainingDays,
		})
	}
	return out, nil
}

func (s *rankingService) MonthlyTotalVolume(
	ctx context.Context, year, month int,
) ([]TotalVolumeDTO, error) {
	from, to, err := calcMonthRange(year, month)
	if err != nil {
		return nil, err
	}

	rows, err := s.repo.MonthlyTotalVolume(ctx, from, to)
	if err != nil {
		return nil, err
	}

	out := make([]TotalVolumeDTO, 0, len(rows))
	for _, r := range rows {
		out = append(out, TotalVolumeDTO{
			UserID:      r.UserID,
			Email:       r.Email,
			TotalVolume: r.TotalVolume,
		})
	}
	return out, nil
}
