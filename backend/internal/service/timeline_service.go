package service

import (
	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
)

type TimelineService interface {
	GetTimeline() ([]repository.TimelineItem, error)
}

type timelineService struct {
	repo repository.TimelineRepository
}

func NewTimelineService(repo repository.TimelineRepository) TimelineService {
	return &timelineService{repo: repo}
}

func (s *timelineService) GetTimeline() ([]repository.TimelineItem, error) {
	return s.repo.FindPublicRecords()
}
