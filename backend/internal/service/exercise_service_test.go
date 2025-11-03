package service

import (
	"context"
	"errors"
	"testing"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/stretchr/testify/require"
)

type fakeExerciseRepo struct {
	listFunc func(ctx context.Context) ([]models.Exercise, error)
}

func (f *fakeExerciseRepo) List(ctx context.Context) ([]models.Exercise, error) {
	return f.listFunc(ctx)
}

func TestExerciseService_List(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		repo        fakeExerciseRepo
		wantDTOs    []ExerciseDTO
		wantErr     bool
		errContains string
	}{
		{
			name: "【正常系】レコードが存在する場合、DTOリストを返すこと",
			repo: fakeExerciseRepo{
				listFunc: func(ctx context.Context) ([]models.Exercise, error) {
					return []models.Exercise{
						{Name: "Bench Press"},
						{Name: "Squat"},
					}, nil
				},
			},
			wantDTOs: []ExerciseDTO{
				{Name: "Bench Press"},
				{Name: "Squat"},
			},
			wantErr: false,
		},
		{
			name: "【正常系】レコードが0件の場合、空スライスを返すこと",
			repo: fakeExerciseRepo{
				listFunc: func(ctx context.Context) ([]models.Exercise, error) {
					return []models.Exercise{}, nil
				},
			},
			wantDTOs: []ExerciseDTO{},
			wantErr:  false,
		},
		{
			name: "【異常系】リポジトリがエラーを返した場合、エラーが伝搬されること",
			repo: fakeExerciseRepo{
				listFunc: func(ctx context.Context) ([]models.Exercise, error) {
					return nil, errors.New("db down")
				},
			},
			wantErr:     true,
			errContains: "種目一覧の取得に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewExerciseService(&tt.repo)
			got, err := svc.List(ctx)

			switch {
			case tt.wantErr:
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errContains)
				require.Nil(t, got)
			default:
				require.NoError(t, err)
				require.Equal(t, tt.wantDTOs, got)
			}
		})
	}
}
