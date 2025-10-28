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

	slog.InfoContext(ctx, "home_summary_fetched",
		"user_id", userID,
		"total_training_days", summary.TotalTrainingDays,
		"latest_weight", summary.LatestWeight,
		"latest_trained_on", summary.LatestTrainedOn,
		"goal_weight", summary.GoalWeight,
		"height", summary.Height,
	)

	return c.JSON(http.StatusOK, summary)
}
