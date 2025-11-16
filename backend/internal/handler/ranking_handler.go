// internal/handler/ranking_handler.go
package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/httpx"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
)

type RankingHandler interface {
	MonthlyGymDays(c echo.Context) error
}

type rankingHandler struct {
	svc   service.RankingService
	cache *service.RankingCache
}

func NewRankingHandler(svc service.RankingService, cache *service.RankingCache) RankingHandler {
	return &rankingHandler{svc: svc, cache: cache}
}

func (h *rankingHandler) MonthlyGymDays(c echo.Context) error {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	ctx := c.Request().Context()

	list, err := h.svc.MonthlyGymDays(ctx, year, month)
	fmt.Println("list", list)
	if err != nil {
		return httpx.Internal("ジム日数ランキングの取得に失敗しました", err)
	}

	h.cache.SetGymDays(list)

	slog.InfoContext(
		ctx, "monthly_gym_days_ranking_fetched",
		"year", year,
		"month", month,
		"count", len(list))

	return c.JSON(http.StatusOK, list)
}
