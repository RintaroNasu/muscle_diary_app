// internal/handler/ranking_handler.go
package handler

import (
	"context"
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

	cached, lastUpdated := h.cache.GetGymDays()

	isFresh := !lastUpdated.IsZero() &&
		lastUpdated.Year() == year &&
		int(lastUpdated.Month()) == month

	var result []service.GymDaysDTO
	var err error

	if !isFresh || len(cached) == 0 {
		result, err = h.svc.MonthlyGymDays(ctx, year, month)
		if err != nil {
			return httpx.Internal("ジム日数ランキングの取得に失敗しました", err)
		}

		h.cache.SetGymDays(result)

		slog.InfoContext(
			ctx, "monthly_gym_days_ranking_fetched_from_db",
			"year", year,
			"month", month,
			"count", len(result),
		)
	} else {
		result = cached

		slog.InfoContext(
			ctx, "monthly_gym_days_ranking_served_from_cache",
			"year", year,
			"month", month,
			"count", len(result),
		)
	}

	go func(year, month int) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		fresh, err := h.svc.MonthlyGymDays(bgCtx, year, month)
		if err != nil {
			slog.ErrorContext(bgCtx, "refresh_monthly_gym_days_failed",
				"year", year,
				"month", month,
				"err", err,
			)
			return
		}

		h.cache.SetGymDays(fresh)

		slog.InfoContext(
			bgCtx, "monthly_gym_days_ranking_refreshed",
			"year", year,
			"month", month,
			"count", len(fresh),
		)
	}(year, month)

	return c.JSON(http.StatusOK, result)
}
