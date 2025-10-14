package handler

import (
	"log/slog"
	"net/http"

	"github.com/RintaroNasu/muscle_diary_app/internal/httpx"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
)

type ExerciseHandler interface {
	List(c echo.Context) error
}

type exerciseHandler struct {
	svc service.ExerciseService
}

func NewExerciseHandler(svc service.ExerciseService) ExerciseHandler {
	return &exerciseHandler{svc: svc}
}

func (h *exerciseHandler) List(c echo.Context) error {
	ctx := c.Request().Context()

	items, err := h.svc.List(c.Request().Context())
	if err != nil {
		return httpx.Internal("システムエラーが発生しました", err)
	}

	slog.InfoContext(ctx, "exercise_list_fetched", "count", len(items))

	return c.JSON(http.StatusOK, items)
}
