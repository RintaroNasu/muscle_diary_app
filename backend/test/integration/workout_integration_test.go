// internal/integration/workout_integration_test.go
package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/RintaroNasu/muscle_diary_app/internal/handler"
	"github.com/RintaroNasu/muscle_diary_app/internal/httpx"
	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newWorkoutIntegrationDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&models.User{},
		&models.Exercise{},
		&models.WorkoutRecord{},
		&models.WorkoutSet{},
	))
	return db
}

func newEchoWithErrHandler() *echo.Echo {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	e.HTTPErrorHandler = httpx.HTTPErrorHandler(logger)
	return e
}

func TestWorkoutIntegration_NormalCases(t *testing.T) {
	e := newEchoWithErrHandler()

	type created struct {
		RecordID uint `json:"record_id"`
	}

	t.Run("【正常系】Create → GetByDate → GetByExercise → Update → Delete", func(t *testing.T) {
		// --- Arrange shared fixtures ---
		db := newWorkoutIntegrationDB(t)

		// ① ユーザー・種目を準備（FK満たす）
		user := models.User{Email: "wk@example.com"}
		require.NoError(t, db.Create(&user).Error)

		ex := models.Exercise{Name: "ベンチプレス"}
		require.NoError(t, db.Create(&ex).Error)

		// ② DI: repo → svc → handler
		repo := repository.NewWorkoutRepository(db)
		svc := service.NewWorkoutService(repo)
		h := handler.NewWorkoutHandler(svc)

		// パラメータ共通化
		const trainedOn = "2025-10-02"

		// ---------------- Create ----------------
		{
			body := `{
				"body_weight": 70.5,
				"exercise_id": %d,
				"trained_on": "` + trainedOn + `",
				"sets": [
					{"set":1,"reps":10,"exercise_weight":50},
					{"set":2,"reps":8,"exercise_weight":55}
				]
			}`
			body = strings.TrimSpace(body)
			body = strings.ReplaceAll(body, "\n", "")
			body = strings.ReplaceAll(body, "\t", "")
			body = strings.Replace(body, "%d", func() string { return strconvUint(ex.ID) }(), 1)

			req := httptest.NewRequest(http.MethodPost, "/workouts", bytes.NewBufferString(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			setUserID(c, user.ID)

			err := h.CreateWorkoutRecord(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, http.StatusCreated, rec.Code)

			var res created
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
			require.NotZero(t, res.RecordID)

			// DB 反映確認（レコード＆セット）
			var got models.WorkoutRecord
			require.NoError(t, db.Preload("Sets").First(&got, res.RecordID).Error)
			require.Equal(t, user.ID, got.UserID)
			require.Equal(t, ex.ID, got.ExerciseID)
			require.Len(t, got.Sets, 2)
		}

		// ---------------- Get by date ----------------
		{
			req := httptest.NewRequest(http.MethodGet, "/workouts?date="+trainedOn, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			setUserID(c, user.ID)

			err := h.GetWorkoutRecordsByDate(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, http.StatusOK, rec.Code)
			// 返却 DTO を軽く確認（exercise_name / sets）
			type setDTO struct {
				Set int `json:"set"`
			}
			type recDTO struct {
				ExerciseName string   `json:"exercise_name"`
				TrainedOn    string   `json:"trained_on"`
				Sets         []setDTO `json:"sets"`
			}
			var arr []recDTO
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &arr))
			require.NotEmpty(t, arr)
			require.Equal(t, "ベンチプレス", arr[0].ExerciseName)
			require.Equal(t, trainedOn, arr[0].TrainedOn)
			require.True(t, len(arr[0].Sets) >= 1)
		}

		// ---------------- Get by exercise ----------------
		{
			req := httptest.NewRequest(http.MethodGet, "/workouts/exercise/"+strconvUint(ex.ID), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("exerciseId")
			c.SetParamValues(strconvUint(ex.ID))
			setUserID(c, user.ID)

			err := h.GetWorkoutRecordsByExercise(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, http.StatusOK, rec.Code)
			// 1件以上返る想定（Createで2セット入れているので >=2）
			var out []map[string]any
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &out))
			require.GreaterOrEqual(t, len(out), 2)
		}

		// ---------------- Update ----------------
		var targetRecord models.WorkoutRecord
		require.NoError(t, db.Where("user_id = ? AND exercise_id = ?", user.ID, ex.ID).First(&targetRecord).Error)

		{
			upBody := `{
				"body_weight": 68.0,
				"exercise_id": %d,
				"trained_on": "2025-10-03",
				"sets": [{"set":1,"reps":6,"exercise_weight":60}]
			}`
			upBody = strings.Replace(upBody, "%d", strconvUint(ex.ID), 1)

			req := httptest.NewRequest(http.MethodPut, "/workouts/"+strconvUint(targetRecord.ID), bytes.NewBufferString(upBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(strconvUint(targetRecord.ID))
			setUserID(c, user.ID)

			err := h.UpdateWorkoutRecord(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, http.StatusOK, rec.Code)

			// DB 更新確認
			var after models.WorkoutRecord
			require.NoError(t, db.Preload("Sets").First(&after, targetRecord.ID).Error)
			require.InDelta(t, 68.0, after.BodyWeight, 1e-6)
			require.Len(t, after.Sets, 1)
			require.Equal(t, 1, after.Sets[0].SetNo)
			require.Equal(t, 6, after.Sets[0].Reps)
		}

		// ---------------- Delete ----------------
		{
			req := httptest.NewRequest(http.MethodDelete, "/workouts/"+strconvUint(targetRecord.ID), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(strconvUint(targetRecord.ID))
			setUserID(c, user.ID)

			err := h.DeleteWorkoutRecord(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, http.StatusOK, rec.Code)

			var cnt int64
			require.NoError(t, db.Model(&models.WorkoutRecord{}).Where("id = ?", targetRecord.ID).Count(&cnt).Error)
			require.EqualValues(t, 0, cnt) // hard-delete(Unscoped) を確認
		}
	})
}

// ---- small util ----
func strconvUint(v uint) string { return strconv.FormatUint(uint64(v), 10) }
