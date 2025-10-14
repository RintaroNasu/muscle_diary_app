package service

import (
	"context"
	"fmt"

	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
)

type ExerciseService interface {
	List(ctx context.Context) ([]ExerciseDTO, error)
}

type ExerciseDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type exerciseService struct {
	repo repository.ExerciseRepository
}

func NewExerciseService(repo repository.ExerciseRepository) ExerciseService {
	return &exerciseService{repo: repo}
}

func (s *exerciseService) List(ctx context.Context) ([]ExerciseDTO, error) {
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("種目一覧の取得に失敗しました: %w", err)
	}

	exercises := make([]ExerciseDTO, 0, len(rows))
	for _, m := range rows {
		exercises = append(exercises, ExerciseDTO{
			ID:   m.ID,
			Name: m.Name,
		})
	}

	return exercises, nil
}
