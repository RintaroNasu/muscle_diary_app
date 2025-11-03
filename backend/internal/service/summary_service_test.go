package service

import (
	"errors"
	"testing"
	"time"

	"github.com/RintaroNasu/muscle_diary_app/utils"
	"github.com/stretchr/testify/require"
)

type fakeSummaryRepo struct {
	countFn  func(userID uint) (int64, error)
	latestFn func(userID uint) (*float64, *time.Time, error)
	basicsFn func(userID uint) (*float64, *float64, error)
}

func (f *fakeSummaryRepo) CountTrainingDays(userID uint) (int64, error) {
	return f.countFn(userID)
}
func (f *fakeSummaryRepo) GetLatestWeight(userID uint) (*float64, *time.Time, error) {
	return f.latestFn(userID)
}
func (f *fakeSummaryRepo) GetProfileBasics(userID uint) (*float64, *float64, error) {
	return f.basicsFn(userID)
}

func TestNewSummaryService(t *testing.T) {
	svc := NewSummaryService(&fakeSummaryRepo{
		countFn:  func(uint) (int64, error) { return 0, nil },
		latestFn: func(uint) (*float64, *time.Time, error) { return nil, nil, nil },
		basicsFn: func(uint) (*float64, *float64, error) { return nil, nil, nil },
	})
	require.NotNil(t, svc)
	_, ok := svc.(SummaryService)
	require.True(t, ok)
}

func TestSummaryService_GetHomeSummary(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		repo        fakeSummaryRepo
		userID      uint
		wantDays    int64
		wantWeight  *float64
		wantTrained *time.Time
		wantHeight  *float64
		wantGoal    *float64
		wantErr     string
	}{
		{
			name: "【正常系】全てのホームサマリ情報が取得できること",
			repo: fakeSummaryRepo{
				countFn: func(uint) (int64, error) { return 12, nil },
				latestFn: func(uint) (*float64, *time.Time, error) {
					return utils.Ptr(62.3), &now, nil
				},
				basicsFn: func(uint) (*float64, *float64, error) {
					return utils.Ptr(175.0), utils.Ptr(65.0), nil
				},
			},
			userID:      1,
			wantDays:    12,
			wantWeight:  utils.Ptr(62.3),
			wantTrained: &now,
			wantHeight:  utils.Ptr(175.0),
			wantGoal:    utils.Ptr(65.0),
		},
		{
			name: "【正常系】最新体重が存在しない場合は nil が返ること",
			repo: fakeSummaryRepo{
				countFn: func(uint) (int64, error) { return 0, nil },
				latestFn: func(uint) (*float64, *time.Time, error) {
					return nil, nil, nil
				},
				basicsFn: func(uint) (*float64, *float64, error) {
					return utils.Ptr(170.0), utils.Ptr(60.0), nil
				},
			},
			userID:      9,
			wantDays:    0,
			wantWeight:  nil,
			wantTrained: nil,
			wantHeight:  utils.Ptr(170.0),
			wantGoal:    utils.Ptr(60.0),
		},
		{
			name: "【異常系】CountTrainingDays のエラーをそのまま返す",
			repo: fakeSummaryRepo{
				countFn:  func(uint) (int64, error) { return 0, errors.New("count failed") },
				latestFn: func(uint) (*float64, *time.Time, error) { return nil, nil, nil },
				basicsFn: func(uint) (*float64, *float64, error) { return nil, nil, nil },
			},
			userID:  1,
			wantErr: "count failed",
		},
		{
			name: "【異常系】GetLatestWeight のエラーをそのまま返す",
			repo: fakeSummaryRepo{
				countFn: func(uint) (int64, error) { return 3, nil },
				latestFn: func(uint) (*float64, *time.Time, error) {
					return nil, nil, errors.New("latest failed")
				},
				basicsFn: func(uint) (*float64, *float64, error) { return nil, nil, nil },
			},
			userID:  1,
			wantErr: "latest failed",
		},
		{
			name: "【異常系】GetProfileBasics のエラーをそのまま返す",
			repo: fakeSummaryRepo{
				countFn:  func(uint) (int64, error) { return 3, nil },
				latestFn: func(uint) (*float64, *time.Time, error) { return nil, nil, nil },
				basicsFn: func(uint) (*float64, *float64, error) { return nil, nil, errors.New("basics failed") },
			},
			userID:  1,
			wantErr: "basics failed",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := NewSummaryService(&tt.repo)

			got, err := svc.GetHomeSummary(tt.userID)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)

			require.Equal(t, tt.wantDays, got.TotalTrainingDays)

			if tt.wantWeight == nil {
				require.Nil(t, got.LatestWeight)
			} else {
				require.NotNil(t, got.LatestWeight)
				require.InDelta(t, *tt.wantWeight, *got.LatestWeight, 1e-6)
			}

			if tt.wantTrained == nil {
				require.Nil(t, got.LatestTrainedOn)
			} else {
				require.NotNil(t, got.LatestTrainedOn)
				require.WithinDuration(t, *tt.wantTrained, *got.LatestTrainedOn, time.Second)
			}

			if tt.wantHeight == nil {
				require.Nil(t, got.Height)
			} else {
				require.NotNil(t, got.Height)
				require.InDelta(t, *tt.wantHeight, *got.Height, 1e-6)
			}
			if tt.wantGoal == nil {
				require.Nil(t, got.GoalWeight)
			} else {
				require.NotNil(t, got.GoalWeight)
				require.InDelta(t, *tt.wantGoal, *got.GoalWeight, 1e-6)
			}
		})
	}
}
