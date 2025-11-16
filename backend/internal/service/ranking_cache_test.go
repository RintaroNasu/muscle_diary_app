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

func TestRankingCache_GetGymDays_ReturnsCopy(t *testing.T) {
	cache := NewRankingCache()

	orig := []GymDaysDTO{
		{UserID: 1, Email: "a@example.com", TotalTrainingDays: 3},
	}
	cache.SetGymDays(orig)

	got1, _ := cache.GetGymDays()
	require.Len(t, got1, 1)
	require.Equal(t, uint(1), got1[0].UserID)

	got1[0].Email = "modified@example.com"

	got2, _ := cache.GetGymDays()

	require.Equal(t, "a@example.com", got2[0].Email)
}
