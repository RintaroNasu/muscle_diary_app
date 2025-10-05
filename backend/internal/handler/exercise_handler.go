package handler

import (
	"fmt"
	"net/http"

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
	items, err := h.svc.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("failed to fetch exercises: %v", err),
		})
	}
	return c.JSON(http.StatusOK, items)
}
