package handler

import (
	"log/slog"
	"net/http"

	"github.com/RintaroNasu/muscle_diary_app/internal/httpx"
	"github.com/RintaroNasu/muscle_diary_app/internal/middleware"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
)

type SummaryHandler interface {
	GetHomeSummary(c echo.Context) error
}

type summaryHandler struct {
	svc service.SummaryService
}

type SummaryResponse struct {
	TotalTrainingDays int      `json:"total_training_days"`
	LatestWeight      *float64 `json:"latest_weight"`
	LatestTrainedOn   string   `json:"latest_trained_on"`
	GoalWeight        *float64 `json:"goal_weight"`
	Height            *float64 `json:"height"`
}

func NewSummaryHandler(svc service.SummaryService) SummaryHandler {
	return &summaryHandler{svc: svc}
}

func (h *summaryHandler) GetHomeSummary(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)

	summary, err := h.svc.GetHomeSummary(userID)
	if err != nil {
		return httpx.Internal("サマリーの取得に失敗しました", err)
	}

	var latestOn string
	if summary.LatestTrainedOn != nil {
		latestOn = summary.LatestTrainedOn.Format("2006-01-02")
	}

	res := SummaryResponse{
		TotalTrainingDays: int(summary.TotalTrainingDays),
		LatestWeight:      summary.LatestWeight,
		LatestTrainedOn:   latestOn,
		GoalWeight:        summary.GoalWeight,
		Height:            summary.Height,
	}

	slog.InfoContext(ctx, "home_summary_fetched",
		"user_id", userID,
		"total_training_days", res.TotalTrainingDays,
		"latest_weight", res.LatestWeight,
		"latest_trained_on", res.LatestTrainedOn,
		"goal_weight", res.GoalWeight,
		"height", res.Height,
	)

	return c.JSON(http.StatusOK, res)
}
