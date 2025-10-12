package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/middleware"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
)

type WorkoutHandler interface {
	CreateWorkoutRecord(c echo.Context) error
	GetWorkoutRecordsByDate(c echo.Context) error
	GetMonthRecordDays(c echo.Context) error
}

type workoutHandler struct {
	svc service.WorkoutService
}

func NewWorkoutHandler(svc service.WorkoutService) WorkoutHandler {
	return &workoutHandler{svc: svc}
}

type CreateWorkoutRecordRequest struct {
	BodyWeight float64             `json:"body_weight"`
	ExerciseID uint                `json:"exercise_id"`
	Sets       []WorkoutSetRequest `json:"sets"`
	TrainedOn  string              `json:"trained_on"`
}

type WorkoutSetRequest struct {
	Set            int     `json:"set"`
	Reps           int     `json:"reps"`
	ExerciseWeight float64 `json:"exercise_weight"`
}

func (h *workoutHandler) CreateWorkoutRecord(c echo.Context) error {
	var req CreateWorkoutRecordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("invalid request body: %v", err),
		})
	}
	userID := middleware.GetUserID(c)

	trainedOn, err := time.Parse("2006-01-02", req.TrainedOn)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("invalid date format: %v", err),
		})
	}

	var sets []service.WorkoutSetData
	for _, setReq := range req.Sets {
		sets = append(sets, service.WorkoutSetData{
			SetNo:          setReq.Set,
			Reps:           setReq.Reps,
			ExerciseWeight: setReq.ExerciseWeight,
		})
	}

	record, err := h.svc.CreateWorkoutRecord(userID, req.BodyWeight, req.ExerciseID, trainedOn, sets)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":   "Workout record created successfully",
		"record_id": record.ID,
	})
}

type workoutSetDTO struct {
	Set            int     `json:"set"`
	Reps           int     `json:"reps"`
	ExerciseWeight float64 `json:"exercise_weight"`
}

type workoutRecordDTO struct {
	ID           uint            `json:"id"`
	ExerciseName string          `json:"exercise_name"`
	BodyWeight   float64         `json:"body_weight"`
	TrainedOn    string          `json:"trained_on"`
	Sets         []workoutSetDTO `json:"sets"`
}

func (h *workoutHandler) GetWorkoutRecordsByDate(c echo.Context) error {
	userID := middleware.GetUserID(c)

	dateStr := c.QueryParam("date")
	if dateStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "date parameter is required (format: YYYY-MM-DD)",
		})
	}

	loc, _ := time.LoadLocation("Asia/Tokyo")
	day, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	fmt.Println("day", day)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("invalid date format: %v", err),
		})
	}

	records, err := h.svc.GetDailyRecords(userID, day)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("failed to get workout records: %v", err),
		})
	}

	out := make([]workoutRecordDTO, 0, len(records))
	for _, r := range records {
		name := r.Exercise.Name
		if name == "" {
			name = "Unknown"
		}
		sets := make([]workoutSetDTO, 0, len(r.Sets))
		for _, s := range r.Sets {
			sets = append(sets, workoutSetDTO{
				Set:            s.SetNo,
				Reps:           s.Reps,
				ExerciseWeight: s.ExerciseWeight,
			})
		}
		out = append(out, workoutRecordDTO{
			ID:           r.ID,
			ExerciseName: name,
			BodyWeight:   r.BodyWeight,
			TrainedOn:    r.TrainedOn.Format("2006-01-02"),
			Sets:         sets,
		})
	}

	return c.JSON(http.StatusOK, out)
}

func (h *workoutHandler) GetMonthRecordDays(c echo.Context) error {
	userID := middleware.GetUserID(c)

	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	if yearStr == "" || monthStr == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "year and month are required (e.g. ?year=2025&month=10)",
			})
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid year"})
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid month"})
	}

	days, err := h.svc.GetMonthRecordDays(userID, year, month)
	if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": fmt.Sprintf("failed to get month record days: %v", err),
			})
	}

	out := make([]string, 0, len(days))
	for _, d := range days {
			out = append(out, d.Format("2006-01-02"))
	}
	return c.JSON(http.StatusOK, out)
}
