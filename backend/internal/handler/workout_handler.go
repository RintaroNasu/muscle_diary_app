package handler

import (
	"net/http"
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/middleware"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
)

type WorkoutHandler interface {
	CreateWorkoutRecord(c echo.Context) error
}

type workoutHandler struct {
	svc service.WorkoutService
}

func NewWorkoutHandler(svc service.WorkoutService) WorkoutHandler {
	return &workoutHandler{svc: svc}
}

type CreateWorkoutRecordRequest struct {
	BodyWeight   float64             `json:"body_weight"`
	ExerciseName string              `json:"exercise_name"`
	Sets         []WorkoutSetRequest `json:"sets"`
	TrainedAt    string              `json:"trained_at"`
}

type WorkoutSetRequest struct {
	Set            int     `json:"set"`
	Reps           int     `json:"reps"`
	ExerciseWeight float64 `json:"exercise_weight"`
}

func (h *workoutHandler) CreateWorkoutRecord(c echo.Context) error {

	var req CreateWorkoutRecordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	userID := middleware.GetUserID(c)

	trainedAt, err := time.Parse(time.RFC3339, req.TrainedAt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid date format"})
	}

	var sets []service.WorkoutSetData
	for _, setReq := range req.Sets {
		sets = append(sets, service.WorkoutSetData{
			SetNo:          setReq.Set,
			Reps:           setReq.Reps,
			ExerciseWeight: setReq.ExerciseWeight,
		})
	}

	record, err := h.svc.CreateWorkoutRecord(userID, req.BodyWeight, req.ExerciseName, trainedAt, sets)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":   "Workout record created successfully",
		"record_id": record.ID,
	})
}
