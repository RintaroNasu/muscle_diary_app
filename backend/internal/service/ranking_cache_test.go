package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRankingCache_SetAndGetGymDays(t *testing.T) {
	cache := NewRankingCache()

	gotData, gotTime := cache.GetGymDays()
	require.Len(t, gotData, 0)
	require.True(t, gotTime.IsZero())

	data := []GymDaysDTO{
		{
			UserID:            1,
			Email:             "test@example.com",
			TotalTrainingDays: 3,
		},
	}

	cache.SetGymDays(data)

	gotData2, gotTime2 := cache.GetGymDays()
	require.Len(t, gotData2, 1)
	require.Equal(t, uint(1), gotData2[0].UserID)
	require.Equal(t, "test@example.com", gotData2[0].Email)
	require.Equal(t, int64(3), gotData2[0].TotalTrainingDays)
	require.False(t, gotTime2.IsZero())
}
