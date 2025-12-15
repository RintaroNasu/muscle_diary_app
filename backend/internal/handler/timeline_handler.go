package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/httpx"
	"github.com/RintaroNasu/muscle_diary_app/internal/middleware"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
)

type TimelineHandler interface {
	GetTimeline(c echo.Context) error
}

type timelineHandler struct {
	svc service.TimelineService
}

func NewTimelineHandler(svc service.TimelineService) TimelineHandler {
	return &timelineHandler{svc: svc}
}

type TimelineItemResponse struct {
	RecordID     uint    `json:"record_id"`
	UserID       uint    `json:"user_id"`
	UserEmail    string  `json:"user_email"`
	ExerciseName string  `json:"exercise_name"`
	BodyWeight   float64 `json:"body_weight"`
	TrainedOn    string  `json:"trained_on"`
	Comment      string  `json:"comment"`
	LikedByMe    bool    `json:"liked_by_me"`
}

func (h *timelineHandler) GetTimeline(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)

	items, err := h.svc.GetTimeline(userID)
	if err != nil {
		return httpx.Internal("システムエラーが発生しました", err)
	}

	loc, _ := time.LoadLocation("Asia/Tokyo")

	var res []TimelineItemResponse
	for _, it := range items {
		res = append(res, TimelineItemResponse{
			RecordID:     it.RecordID,
			UserID:       it.UserID,
			UserEmail:    it.UserEmail,
			ExerciseName: it.ExerciseName,
			BodyWeight:   it.BodyWeight,
			TrainedOn:    it.TrainedOn.In(loc).Format("2006-01-02"),
			Comment:      it.Comment,
			LikedByMe:    it.LikedByMe,
		})
	}

	slog.InfoContext(ctx, "timeline_fetched",
		"count", len(res),
	)

	return c.JSON(http.StatusOK, res)
}
