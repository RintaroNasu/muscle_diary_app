package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/httpx"
	"github.com/RintaroNasu/muscle_diary_app/internal/middleware"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
)

type WorkoutHandler interface {
	CreateWorkoutRecord(c echo.Context) error
	GetWorkoutRecordsByDate(c echo.Context) error
	GetMonthRecordDays(c echo.Context) error
	UpdateWorkoutRecord(c echo.Context) error
	DeleteWorkoutRecord(c echo.Context) error
	GetWorkoutRecordsByExercise(c echo.Context) error
}

type workoutHandler struct {
	svc service.WorkoutService
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

type ExerciseSingleSetResponse struct {
	RecordID       uint    `json:"record_id"`
	TrainedOn      string  `json:"trained_on"`
	Set            int     `json:"set"`
	Reps           int     `json:"reps"`
	ExerciseWeight float64 `json:"exercise_weight"`
	BodyWeight     float64 `json:"body_weight"`
}

func NewWorkoutHandler(svc service.WorkoutService) WorkoutHandler {
	return &workoutHandler{svc: svc}
}

func (h *workoutHandler) CreateWorkoutRecord(c echo.Context) error {
	ctx := c.Request().Context()
	var req CreateWorkoutRecordRequest

	if err := c.Bind(&req); err != nil {
		return httpx.BadRequest("InvalidBody", "リクエストの形式が不正です", err)
	}
	userID := middleware.GetUserID(c)

	loc, _ := time.LoadLocation("Asia/Tokyo")
	trainedOn, err := time.ParseInLocation("2006-01-02", req.TrainedOn, loc)
	if err != nil {
		return httpx.BadRequest("InvalidDate", "日付の形式が不正です", err)
	}
	trainedOn = time.Date(trainedOn.Year(), trainedOn.Month(), trainedOn.Day(), 0, 0, 0, 0, time.UTC)

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
		switch {
		case errors.Is(err, service.ErrNoSets),
			errors.Is(err, service.ErrInvalidSetValue):
			return httpx.BadRequest("ValidationError", "セット内容が不正です", err)
		case errors.Is(err, service.ErrExerciseNotFound):
			return httpx.NotFound("ExerciseNotFound", "指定の種目が見つかりません", err)
		default:
			return httpx.Internal("システムエラーが発生しました", err)
		}
	}

	slog.InfoContext(ctx, "workout_created",
		"record_id", record.ID,
		"exercise_id", req.ExerciseID,
	)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":   "Workout record created successfully",
		"record_id": record.ID,
	})
}

func (h *workoutHandler) GetWorkoutRecordsByDate(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)

	dateStr := c.QueryParam("date")
	if dateStr == "" {
		return httpx.BadRequest("InvalidQuery", "date は必須です（YYYY-MM-DD）", nil)
	}

	loc, _ := time.LoadLocation("Asia/Tokyo")
	day, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		return httpx.BadRequest("InvalidDate", "date の形式が不正です（YYYY-MM-DD）", err)
	}
	day = time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)

	records, err := h.svc.GetDailyRecords(userID, day)
	if err != nil {
		return httpx.Internal("システムエラーが発生しました", err)
	}

	slog.InfoContext(ctx, "workout_daily_fetched",
		"date", day.Format("2006-01-02"),
		"count", len(records),
	)

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
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)

	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	if yearStr == "" || monthStr == "" {
		return httpx.BadRequest("InvalidQuery", "year と month は必須です（例: ?year=2025&month=10）", nil)
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return httpx.BadRequest("InvalidYear", "year の形式が不正です（整数）", err)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return httpx.BadRequest("InvalidMonth", "month は 1〜12 の整数です", err)
	}

	days, err := h.svc.GetMonthRecordDays(userID, year, month)
	if err != nil {
		return httpx.Internal("システムエラーが発生しました", err)
	}

	slog.InfoContext(ctx, "workout_month_days_fetched",
		"year", year,
		"month", month,
		"count", len(days),
	)

	out := make([]string, 0, len(days))
	for _, d := range days {
		out = append(out, d.Format("2006-01-02"))
	}
	return c.JSON(http.StatusOK, out)
}

func (h *workoutHandler) UpdateWorkoutRecord(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)

	recordIDStr := c.Param("id")
	recordID, err := strconv.ParseUint(recordIDStr, 10, 32)
	if err != nil {
		return httpx.BadRequest("InvalidID", "レコードIDが不正です", err)
	}

	var req CreateWorkoutRecordRequest
	if err := c.Bind(&req); err != nil {
		return httpx.BadRequest("InvalidBody", "リクエストの形式が不正です", err)
	}

	loc, _ := time.LoadLocation("Asia/Tokyo")
	trainedOn, err := time.ParseInLocation("2006-01-02", req.TrainedOn, loc)
	if err != nil {
		return httpx.BadRequest("InvalidDate", "日付の形式が不正です", err)
	}
	trainedOn = time.Date(trainedOn.Year(), trainedOn.Month(), trainedOn.Day(), 0, 0, 0, 0, time.UTC)

	var sets []service.WorkoutSetData
	for _, setReq := range req.Sets {
		sets = append(sets, service.WorkoutSetData{
			SetNo:          setReq.Set,
			Reps:           setReq.Reps,
			ExerciseWeight: setReq.ExerciseWeight,
		})
	}

	record, err := h.svc.UpdateWorkoutRecord(userID, uint(recordID), req.BodyWeight, req.ExerciseID, trainedOn, sets)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNoSets),
			errors.Is(err, service.ErrInvalidSetValue):
			return httpx.BadRequest("ValidationError", "セット内容が不正です", err)
		case errors.Is(err, service.ErrExerciseNotFound):
			return httpx.NotFound("ExerciseNotFound", "指定の種目が見つかりません", err)
		case errors.Is(err, service.ErrRecordNotFound):
			return httpx.NotFound("RecordNotFound", "指定の記録が見つかりません", err)
		default:
			return httpx.Internal("システムエラーが発生しました", err)
		}
	}

	slog.InfoContext(ctx, "workout_updated",
		"record_id", record.ID,
		"exercise_id", req.ExerciseID,
	)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "Workout record updated successfully",
		"record_id": record.ID,
	})
}

func (h *workoutHandler) DeleteWorkoutRecord(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)

	recordIDStr := c.Param("id")
	recordID, err := strconv.ParseUint(recordIDStr, 10, 32)
	if err != nil {
		return httpx.BadRequest("InvalidID", "レコードIDが不正です", err)
	}

	err = h.svc.DeleteWorkoutRecord(userID, uint(recordID))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			return httpx.NotFound("RecordNotFound", "指定の記録が見つかりません", err)
		default:
			return httpx.Internal("システムエラーが発生しました", err)
		}
	}

	slog.InfoContext(ctx, "workout_deleted",
		"record_id", recordID,
	)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Workout record deleted successfully",
	})
}

func (h *workoutHandler) GetWorkoutRecordsByExercise(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)
	exerciseIDStr := c.Param("exerciseId")
	if exerciseIDStr == "" {
		return httpx.BadRequest("MissingExerciseID", "種目IDが指定されていません", nil)
	}

	eid, err := strconv.ParseUint(exerciseIDStr, 10, 32)
	if err != nil {
		return httpx.BadRequest("InvalidExerciseID", "種目IDが不正です", err)
	}

	rows, err := h.svc.GetWorkoutRecordsByExercise(userID, uint(eid))
	if err != nil {
		return httpx.Internal("システムエラーが発生しました", err)
	}

	loc, _ := time.LoadLocation("Asia/Tokyo")
	out := make([]ExerciseSingleSetResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, ExerciseSingleSetResponse{
			RecordID:       r.RecordID,
			TrainedOn:      r.TrainedOn.In(loc).Format("2006-01-02"),
			Set:            r.SetNo,
			Reps:           r.Reps,
			ExerciseWeight: r.ExerciseWeight,
			BodyWeight:     r.BodyWeight,
		})
	}

	slog.InfoContext(ctx, "workout_exercise_sets_fetched",
		"exercise_id", eid,
		"count", len(out),
	)

	return c.JSON(http.StatusOK, out)
}
