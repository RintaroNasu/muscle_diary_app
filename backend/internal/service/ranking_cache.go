package service

import (
	"sync"
	"time"
)

type RankingCache struct {
	mu          sync.RWMutex
	gymDays     []GymDaysDTO
	lastUpdated time.Time
}

func NewRankingCache() *RankingCache {
	return &RankingCache{}
}

func (c *RankingCache) GetGymDays() ([]GymDaysDTO, time.Time) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.gymDays) == 0 {
		return nil, c.lastUpdated
	}

	out := make([]GymDaysDTO, len(c.gymDays))
	copy(out, c.gymDays)
	return out, c.lastUpdated
}

func (c *RankingCache) SetGymDays(data []GymDaysDTO) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(data) == 0 {
		c.gymDays = nil
		c.lastUpdated = time.Time{}
		return
	}

	c.gymDays = make([]GymDaysDTO, len(data))
	copy(c.gymDays, data)
	c.lastUpdated = time.Now()
}
