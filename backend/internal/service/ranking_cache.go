package service

import "time"

type RankingCache struct {
	gymDays     []GymDaysDTO
	lastUpdated time.Time
}

func NewRankingCache() *RankingCache {
	return &RankingCache{}
}

// ※ MR1時点では単純にそのまま返す（後でMR2でコピー＋Mutexに変える）
func (c *RankingCache) GetGymDays() ([]GymDaysDTO, time.Time) {
	return c.gymDays, c.lastUpdated
}

func (c *RankingCache) SetGymDays(data []GymDaysDTO) {
	c.gymDays = data
	c.lastUpdated = time.Now()
}
