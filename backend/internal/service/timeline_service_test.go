package service

import (
	"errors"
	"testing"
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
	"github.com/stretchr/testify/require"
)

type fakeTimelineRepo struct {
	findFn func(userID uint) ([]repository.TimelineItem, error)
}

func (f *fakeTimelineRepo) FindPublicRecords(userID uint) ([]repository.TimelineItem, error) {
	return f.findFn(userID)
}
func TestNewTimelineService(t *testing.T) {
	repo := &fakeTimelineRepo{
		findFn: func(userID uint) ([]repository.TimelineItem, error) {
			return nil, nil
		},
	}

	svc := NewTimelineService(repo)
	require.NotNil(t, svc)

	_, ok := svc.(TimelineService)
	require.True(t, ok, "should implement TimelineService")
}
func TestTimelineService_GetTimeline(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		repo      fakeTimelineRepo
		wantLen   int
		wantFirst *TimelineItem
		wantErr   string
	}{
		{
			name: "【正常系】公開記録が1件以上返ってくる場合に値をマッピングできること",
			repo: fakeTimelineRepo{
				findFn: func(userID uint) ([]repository.TimelineItem, error) {
					return []repository.TimelineItem{
						{
							RecordID:     1,
							UserID:       10,
							UserEmail:    "test@example.com",
							ExerciseName: "Bench Press",
							BodyWeight:   70.5,
							TrainedOn:    now,
							Comment:      "がんばった",
							LikedByMe:    false,
						},
					}, nil
				},
			},
			wantLen: 1,
			wantFirst: &TimelineItem{
				RecordID:     1,
				UserID:       10,
				UserEmail:    "test@example.com",
				ExerciseName: "Bench Press",
				BodyWeight:   70.5,
				TrainedOn:    now,
				Comment:      "がんばった",
			},
		},
		{
			name: "【正常系】公開記録が0件の場合は空スライスを返すこと",
			repo: fakeTimelineRepo{
				findFn: func(userID uint) ([]repository.TimelineItem, error) {
					return []repository.TimelineItem{}, nil
				},
			},
			wantLen:   0,
			wantFirst: nil,
		},
		{
			name: "【異常系】リポジトリのエラーをそのまま返すこと",
			repo: fakeTimelineRepo{
				findFn: func(userID uint) ([]repository.TimelineItem, error) {
					return nil, errors.New("db failed")
				},
			},
			wantErr: "db failed",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := NewTimelineService(&tt.repo)

			got, err := svc.GetTimeline(1)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.Len(t, got, tt.wantLen)

			if tt.wantFirst != nil {
				require.Equal(t, tt.wantFirst.RecordID, got[0].RecordID)
				require.Equal(t, tt.wantFirst.UserID, got[0].UserID)
				require.Equal(t, tt.wantFirst.UserEmail, got[0].UserEmail)
				require.Equal(t, tt.wantFirst.ExerciseName, got[0].ExerciseName)
				require.InDelta(t, tt.wantFirst.BodyWeight, got[0].BodyWeight, 1e-6)
				require.WithinDuration(t, tt.wantFirst.TrainedOn, got[0].TrainedOn, time.Second)
				require.Equal(t, tt.wantFirst.Comment, got[0].Comment)
			}
		})
	}
}
