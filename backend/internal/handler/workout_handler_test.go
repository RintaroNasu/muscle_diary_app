package handler

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/httpx"
	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newEchoForTest() *echo.Echo {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	e.HTTPErrorHandler = httpx.HTTPErrorHandler(logger)
	return e
}

type mockWorkoutService struct {
	CreateWorkoutRecordFunc       func(userID uint, bodyWeight float64, exerciseID uint, trainedOn time.Time, sets []service.WorkoutSetData, isPublic bool, comment string) (*models.WorkoutRecord, error)
	GetDailyRecordsFunc           func(userID uint, day time.Time) ([]models.WorkoutRecord, error)
	GetMonthRecordDaysFunc        func(userID uint, year int, month int) ([]time.Time, error)
	UpdateWorkoutRecordFunc       func(userID uint, recordID uint, bodyWeight float64, exerciseID uint, trainedOn time.Time, sets []service.WorkoutSetData) (*models.WorkoutRecord, error)
	DeleteWorkoutRecordFunc       func(userID uint, recordID uint) error
	GetWorkoutRecordsByExerciseFn func(userID uint, exerciseID uint) ([]service.FlatSet, error)
}

func (m *mockWorkoutService) CreateWorkoutRecord(a uint, b float64, c uint, d time.Time, e []service.WorkoutSetData, f bool, g string) (*models.WorkoutRecord, error) {
	return m.CreateWorkoutRecordFunc(a, b, c, d, e, f, g)
}
func (m *mockWorkoutService) GetDailyRecords(a uint, b time.Time) ([]models.WorkoutRecord, error) {
	return m.GetDailyRecordsFunc(a, b)
}
func (m *mockWorkoutService) GetMonthRecordDays(a uint, y int, mo int) ([]time.Time, error) {
	return m.GetMonthRecordDaysFunc(a, y, mo)
}
func (m *mockWorkoutService) UpdateWorkoutRecord(a uint, id uint, bw float64, ex uint, t time.Time, sets []service.WorkoutSetData) (*models.WorkoutRecord, error) {
	return m.UpdateWorkoutRecordFunc(a, id, bw, ex, t, sets)
}
func (m *mockWorkoutService) DeleteWorkoutRecord(a uint, id uint) error {
	return m.DeleteWorkoutRecordFunc(a, id)
}
func (m *mockWorkoutService) GetWorkoutRecordsByExercise(a uint, ex uint) ([]service.FlatSet, error) {
	return m.GetWorkoutRecordsByExerciseFn(a, ex)
}

func TestWorkoutHandler_CreateWorkoutRecord(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		mock         *mockWorkoutService
		wantCode     int
		wantContains string
	}{
		{
			name: "【正常系】レコードとセットを作成できること",
			body: `{"body_weight":70.5,"exercise_id":2,"trained_on":"2025-10-01","sets":[{"set":1,"reps":10,"exercise_weight":50}]}`,
			mock: &mockWorkoutService{
				CreateWorkoutRecordFunc: func(userID uint, bodyWeight float64, exerciseID uint, trainedOn time.Time, sets []service.WorkoutSetData, isPublic bool, comment string) (*models.WorkoutRecord, error) {
					return &models.WorkoutRecord{Model: gorm.Model{ID: 123}}, nil
				},
			},
			wantCode:     http.StatusCreated,
			wantContains: `"record_id":123`,
		},
		{
			name: "【異常系】リクエストの形式が不正な場合は InvalidBody エラーを返すこと",
			body: `{"trained_on":123,"sets":[]}`,
			mock: &mockWorkoutService{
				CreateWorkoutRecordFunc: func(uint, float64, uint, time.Time, []service.WorkoutSetData, bool, string) (*models.WorkoutRecord, error) {
					return nil, nil
				},
			},
			wantCode:     http.StatusBadRequest,
			wantContains: `"code":"InvalidBody"`,
		},
		{
			name: "【異常系】日付の形式が不正な場合は InvalidDate エラーを返すこと",
			body: `{"body_weight":70,"exercise_id":1,"trained_on":"2025/10/01","sets":[{"set":1,"reps":10,"exercise_weight":50}]}`,
			mock: &mockWorkoutService{
				CreateWorkoutRecordFunc: func(uint, float64, uint, time.Time, []service.WorkoutSetData, bool, string) (*models.WorkoutRecord, error) {
					return nil, nil
				},
			},
			wantCode:     http.StatusBadRequest,
			wantContains: `"code":"InvalidDate"`,
		},
		{
			name: "【異常系】セットが空の場合は ErrNoSets を返すこと",
			body: `{"body_weight":70,"exercise_id":1,"trained_on":"2025-10-01","sets":[]}`,
			mock: &mockWorkoutService{
				CreateWorkoutRecordFunc: func(uint, float64, uint, time.Time, []service.WorkoutSetData, bool, string) (*models.WorkoutRecord, error) {
					return nil, service.ErrNoSets
				},
			},
			wantCode:     http.StatusBadRequest,
			wantContains: `"code":"ValidationError"`,
		},
		{
			name: "【異常系】指定の種目が見つからない場合は ExerciseNotFound エラーを返すこと",
			body: `{"body_weight":70,"exercise_id":999,"trained_on":"2025-10-01","sets":[{"set":1,"reps":10,"exercise_weight":50}]}`,
			mock: &mockWorkoutService{
				CreateWorkoutRecordFunc: func(uint, float64, uint, time.Time, []service.WorkoutSetData, bool, string) (*models.WorkoutRecord, error) {
					return nil, service.ErrExerciseNotFound
				},
			},
			wantCode:     http.StatusNotFound,
			wantContains: `"code":"ExerciseNotFound"`,
		},
		{
			name: "【異常系】システムエラーが発生した場合は InternalError エラーを返すこと",
			body: `{"body_weight":70,"exercise_id":1,"trained_on":"2025-10-01","sets":[{"set":1,"reps":10,"exercise_weight":50}]}`,
			mock: &mockWorkoutService{
				CreateWorkoutRecordFunc: func(uint, float64, uint, time.Time, []service.WorkoutSetData, bool, string) (*models.WorkoutRecord, error) {
					return nil, errors.New("db down")
				},
			},
			wantCode:     http.StatusInternalServerError,
			wantContains: `"code":"InternalError"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEchoForTest()
			h := NewWorkoutHandler(tt.mock)

			req := httptest.NewRequest(http.MethodPost, "/workouts", bytes.NewBufferString(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user_id", uint(1))

			err := h.CreateWorkoutRecord(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, tt.wantCode, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantContains)
		})
	}
}

func TestWorkoutHandler_GetWorkoutRecordsByDate(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	day := "2025-10-02"
	parsed, _ := time.ParseInLocation("2006-01-02", day, loc)

	tests := []struct {
		name         string
		query        string
		mock         *mockWorkoutService
		wantCode     int
		wantContains string
	}{
		{
			name:  "【正常系】日別レコードを取得できること",
			query: "date=" + day,
			mock: &mockWorkoutService{
				GetDailyRecordsFunc: func(userID uint, d time.Time) ([]models.WorkoutRecord, error) {
					return []models.WorkoutRecord{
						{
							Model:      gorm.Model{ID: 1},
							BodyWeight: 70,
							TrainedOn:  parsed,
							Exercise:   models.Exercise{Name: "Bench"},
							Sets:       []models.WorkoutSet{{SetNo: 1, Reps: 10, ExerciseWeight: 50}},
						},
						{
							Model:      gorm.Model{ID: 2},
							BodyWeight: 71,
							TrainedOn:  parsed,
							Exercise:   models.Exercise{Name: "Squat"},
							Sets:       []models.WorkoutSet{{SetNo: 1, Reps: 8, ExerciseWeight: 80}},
						},
					}, nil
				},
			},
			wantCode:     http.StatusOK,
			wantContains: `"exercise_name":"Bench"`,
		},
		{
			name:         "【異常系】date が必須であること",
			query:        "",
			mock:         &mockWorkoutService{},
			wantCode:     http.StatusBadRequest,
			wantContains: `"code":"InvalidQuery"`,
		},
		{
			name:         "【異常系】日付の形式が不正な場合は InvalidDate エラーを返すこと",
			query:        "date=2025/10/02",
			mock:         &mockWorkoutService{},
			wantCode:     http.StatusBadRequest,
			wantContains: `"code":"InvalidDate"`,
		},
		{
			name:  "【異常系】システムエラーが発生した場合は InternalError エラーを返すこと",
			query: "date=2025-10-02",
			mock: &mockWorkoutService{
				GetDailyRecordsFunc: func(uint, time.Time) ([]models.WorkoutRecord, error) {
					return nil, errors.New("x")
				},
			},
			wantCode:     http.StatusInternalServerError,
			wantContains: `"code":"InternalError"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEchoForTest()
			h := NewWorkoutHandler(tt.mock)

			req := httptest.NewRequest(http.MethodGet, "/workouts?"+tt.query, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user_id", uint(1))

			err := h.GetWorkoutRecordsByDate(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, tt.wantCode, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantContains)
		})
	}
}

func TestWorkoutHandler_GetMonthRecordDays(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		mock         *mockWorkoutService
		wantCode     int
		wantContains string
	}{
		{
			name:  "【正常系】月内の日付を取得できること",
			query: "year=2025&month=10",
			mock: &mockWorkoutService{
				GetMonthRecordDaysFunc: func(uint, int, int) ([]time.Time, error) {
					return []time.Time{
						time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
						time.Date(2025, 10, 3, 0, 0, 0, 0, time.UTC),
					}, nil
				},
			},
			wantCode:     http.StatusOK,
			wantContains: `"2025-10-01"`,
		},
		{
			name:         "【異常系】year/month が必須であること",
			query:        "",
			mock:         &mockWorkoutService{},
			wantCode:     http.StatusBadRequest,
			wantContains: `"code":"InvalidQuery"`,
		},
		{
			name:         "【異常系】year の形式が不正な場合は InvalidYear エラーを返すこと",
			query:        "year=yy&month=10",
			mock:         &mockWorkoutService{},
			wantCode:     http.StatusBadRequest,
			wantContains: `"code":"InvalidYear"`,
		},
		{
			name:         "【異常系】month の形式が不正な場合は InvalidMonth エラーを返すこと",
			query:        "year=2025&month=13",
			mock:         &mockWorkoutService{},
			wantCode:     http.StatusBadRequest,
			wantContains: `"code":"InvalidMonth"`,
		},
		{
			name:  "【異常系】システムエラーが発生した場合は InternalError エラーを返すこと",
			query: "year=2025&month=10",
			mock: &mockWorkoutService{
				GetMonthRecordDaysFunc: func(uint, int, int) ([]time.Time, error) {
					return nil, errors.New("db")
				},
			},
			wantCode:     http.StatusInternalServerError,
			wantContains: `"code":"InternalError"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEchoForTest()
			h := NewWorkoutHandler(tt.mock)

			req := httptest.NewRequest(http.MethodGet, "/workouts/days?"+tt.query, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user_id", uint(1))

			err := h.GetMonthRecordDays(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, tt.wantCode, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantContains)
		})
	}
}

func TestWorkoutHandler_UpdateWorkoutRecord(t *testing.T) {
	tests := []struct {
		name         string
		pathID       string
		body         string
		mock         *mockWorkoutService
		wantCode     int
		wantContains string
	}{
		{
			name:   "【正常系】レコードとセットを更新できること",
			pathID: "777",
			body:   `{"body_weight":68,"exercise_id":4,"trained_on":"2025-10-05","sets":[{"set":1,"reps":8,"exercise_weight":60}]}`,
			mock: &mockWorkoutService{
				UpdateWorkoutRecordFunc: func(uint, uint, float64, uint, time.Time, []service.WorkoutSetData) (*models.WorkoutRecord, error) {
					return &models.WorkoutRecord{Model: gorm.Model{ID: 777}}, nil
				},
			},
			wantCode:     http.StatusOK,
			wantContains: `"record_id":777`,
		},
		{
			name:         "【異常系】ID の形式が不正な場合は InvalidID エラーを返すこと",
			pathID:       "abc",
			body:         `{}`,
			mock:         &mockWorkoutService{},
			wantCode:     http.StatusBadRequest,
			wantContains: `"code":"InvalidID"`,
		},
		{
			name:         "【異常系】リクエストの形式が不正な場合は InvalidBody エラーを返すこと",
			pathID:       "1",
			body:         `{"trained_on":123}`,
			mock:         &mockWorkoutService{},
			wantCode:     http.StatusBadRequest,
			wantContains: `"code":"InvalidBody"`,
		},
		{
			name:         "【異常系】日付の形式が不正な場合は InvalidDate エラーを返すこと",
			pathID:       "1",
			body:         `{"trained_on":"2025/10/05","sets":[{"set":1,"reps":8,"exercise_weight":60}]}`,
			mock:         &mockWorkoutService{},
			wantCode:     http.StatusBadRequest,
			wantContains: `"code":"InvalidDate"`,
		},
		{
			name:   "【異常系】セットが空の場合は ErrNoSets を返すこと",
			pathID: "1",
			body:   `{"trained_on":"2025-10-05","sets":[]}`,
			mock: &mockWorkoutService{
				UpdateWorkoutRecordFunc: func(uint, uint, float64, uint, time.Time, []service.WorkoutSetData) (*models.WorkoutRecord, error) {
					return nil, service.ErrNoSets
				},
			},
			wantCode:     http.StatusBadRequest,
			wantContains: `"code":"ValidationError"`,
		},
		{
			name:   "【異常系】指定の種目が見つからない場合は ExerciseNotFound エラーを返すこと",
			pathID: "1",
			body:   `{"trained_on":"2025-10-05","sets":[{"set":1,"reps":8,"exercise_weight":60}]}`,
			mock: &mockWorkoutService{
				UpdateWorkoutRecordFunc: func(uint, uint, float64, uint, time.Time, []service.WorkoutSetData) (*models.WorkoutRecord, error) {
					return nil, service.ErrExerciseNotFound
				},
			},
			wantCode:     http.StatusNotFound,
			wantContains: `"code":"ExerciseNotFound"`,
		},
		{
			name:   "【異常系】指定の記録が見つからない場合は RecordNotFound エラーを返すこと",
			pathID: "1",
			body:   `{"trained_on":"2025-10-05","sets":[{"set":1,"reps":8,"exercise_weight":60}]}`,
			mock: &mockWorkoutService{
				UpdateWorkoutRecordFunc: func(uint, uint, float64, uint, time.Time, []service.WorkoutSetData) (*models.WorkoutRecord, error) {
					return nil, service.ErrRecordNotFound
				},
			},
			wantCode:     http.StatusNotFound,
			wantContains: `"code":"RecordNotFound"`,
		},
		{
			name:   "【異常系】システムエラーが発生した場合は InternalError エラーを返すこと",
			pathID: "1",
			body:   `{"trained_on":"2025-10-05","sets":[{"set":1,"reps":8,"exercise_weight":60}]}`,
			mock: &mockWorkoutService{
				UpdateWorkoutRecordFunc: func(uint, uint, float64, uint, time.Time, []service.WorkoutSetData) (*models.WorkoutRecord, error) {
					return nil, errors.New("x")
				},
			},
			wantCode:     http.StatusInternalServerError,
			wantContains: `"code":"InternalError"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEchoForTest()
			h := NewWorkoutHandler(tt.mock)

			req := httptest.NewRequest(http.MethodPut, "/workouts/"+tt.pathID, bytes.NewBufferString(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.pathID)
			c.Set("user_id", uint(1))

			err := h.UpdateWorkoutRecord(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, tt.wantCode, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantContains)
		})
	}
}

func TestWorkoutHandler_DeleteWorkoutRecord(t *testing.T) {
	tests := []struct {
		name         string
		pathID       string
		mock         *mockWorkoutService
		wantCode     int
		wantContains string
	}{
		{
			name:   "【正常系】レコードを削除できること",
			pathID: "10",
			mock: &mockWorkoutService{
				DeleteWorkoutRecordFunc: func(uint, uint) error { return nil },
			},
			wantCode:     http.StatusOK,
			wantContains: `"message"`,
		},
		{
			name:         "【異常系】ID の形式が不正な場合は InvalidID エラーを返すこと",
			pathID:       "abc",
			mock:         &mockWorkoutService{},
			wantCode:     http.StatusBadRequest,
			wantContains: `"code":"InvalidID"`,
		},
		{
			name:   "【異常系】指定の記録が見つからない場合は RecordNotFound エラーを返すこと",
			pathID: "1",
			mock: &mockWorkoutService{
				DeleteWorkoutRecordFunc: func(uint, uint) error { return service.ErrRecordNotFound },
			},
			wantCode:     http.StatusNotFound,
			wantContains: `"code":"RecordNotFound"`,
		},
		{
			name:   "【異常系】システムエラーが発生した場合は InternalError エラーを返すこと",
			pathID: "1",
			mock: &mockWorkoutService{
				DeleteWorkoutRecordFunc: func(uint, uint) error { return errors.New("x") },
			},
			wantCode:     http.StatusInternalServerError,
			wantContains: `"code":"InternalError"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEchoForTest()
			h := NewWorkoutHandler(tt.mock)

			req := httptest.NewRequest(http.MethodDelete, "/workouts/"+tt.pathID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.pathID)
			c.Set("user_id", uint(1))

			err := h.DeleteWorkoutRecord(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, tt.wantCode, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantContains)
		})
	}
}

func TestWorkoutHandler_GetWorkoutRecordsByExercise(t *testing.T) {
	tests := []struct {
		name         string
		pathVal      string
		mock         *mockWorkoutService
		wantCode     int
		wantContains string
	}{
		{
			name:    "【正常系】種目のレコードを取得できること",
			pathVal: "2",
			mock: &mockWorkoutService{
				GetWorkoutRecordsByExerciseFn: func(uint, uint) ([]service.FlatSet, error) {
					return []service.FlatSet{
						{RecordID: 1, TrainedOn: time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC), SetNo: 1, Reps: 10, ExerciseWeight: 50, BodyWeight: 70},
						{RecordID: 1, TrainedOn: time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC), SetNo: 2, Reps: 8, ExerciseWeight: 55, BodyWeight: 70},
					}, nil
				},
			},
			wantCode:     http.StatusOK,
			wantContains: `"record_id":1`,
		},
		{
			name:         "【異常系】種目ID が指定されていない場合は MissingExerciseID エラーを返すこと",
			pathVal:      "",
			mock:         &mockWorkoutService{},
			wantCode:     http.StatusBadRequest,
			wantContains: `"code":"MissingExerciseID"`,
		},
		{
			name:         "【異常系】種目ID の形式が不正な場合は InvalidExerciseID エラーを返すこと",
			pathVal:      "abc",
			mock:         &mockWorkoutService{},
			wantCode:     http.StatusBadRequest,
			wantContains: `"code":"InvalidExerciseID"`,
		},
		{
			name:    "【異常系】システムエラーが発生した場合は InternalError エラーを返すこと",
			pathVal: "2",
			mock: &mockWorkoutService{
				GetWorkoutRecordsByExerciseFn: func(uint, uint) ([]service.FlatSet, error) {
					return nil, errors.New("x")
				},
			},
			wantCode:     http.StatusInternalServerError,
			wantContains: `"code":"InternalError"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newEchoForTest()
			h := NewWorkoutHandler(tt.mock)

			url := "/workouts/exercise"
			if tt.pathVal != "" {
				url += "/" + tt.pathVal
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if tt.pathVal != "" {
				c.SetParamNames("exerciseId")
				c.SetParamValues(tt.pathVal)
			}
			c.Set("user_id", uint(1))

			err := h.GetWorkoutRecordsByExercise(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, tt.wantCode, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantContains)
		})
	}
}
