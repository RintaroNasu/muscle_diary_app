package service

import (
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
)

type TimelineService interface {
	GetTimeline(userID uint) ([]TimelineItem, error)
}

type TimelineItem struct {
	RecordID     uint
	UserID       uint
	UserEmail    string
	ExerciseName string
	BodyWeight   float64
	TrainedOn    time.Time
	Comment      string
	LikedByMe    bool
}

type timelineService struct {
	repo repository.TimelineRepository
}

func NewTimelineService(repo repository.TimelineRepository) TimelineService {
	return &timelineService{repo: repo}
}

func (s *timelineService) GetTimeline(userID uint) ([]TimelineItem, error) {
	rows, err := s.repo.FindPublicRecords(userID)
	if err != nil {
		return nil, err
	}

	out := make([]TimelineItem, 0, len(rows))
	for _, it := range rows {
		out = append(out, TimelineItem{
			RecordID:     it.RecordID,
			UserID:       it.UserID,
			UserEmail:    it.UserEmail,
			ExerciseName: it.ExerciseName,
			BodyWeight:   it.BodyWeight,
			TrainedOn:    it.TrainedOn,
			Comment:      it.Comment,
			LikedByMe:    it.LikedByMe,
		})
	}
	return out, nil
}
