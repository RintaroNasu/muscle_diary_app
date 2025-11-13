package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/handler"
	"github.com/RintaroNasu/muscle_diary_app/internal/httpx"
	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/RintaroNasu/muscle_diary_app/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SummaryResponseDTO struct {
	TotalTrainingDays int      `json:"total_training_days"`
	LatestWeight      *float64 `json:"latest_weight"`
	LatestTrainedOn   string   `json:"latest_trained_on"`
	GoalWeight        *float64 `json:"goal_weight"`
	Height            *float64 `json:"height"`
}

func newSummaryIntegrationDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}, &models.Exercise{}, &models.WorkoutRecord{}))
	return db
}

func TestSummaryIntegration_GetHomeSummary(t *testing.T) {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	e.HTTPErrorHandler = httpx.HTTPErrorHandler(logger)

	t.Run("【正常系】DBのデータを通してサマリーが取得できること", func(t *testing.T) {
		db := newSummaryIntegrationDB(t)

		u := models.User{
			Email:      "user@example.com",
			Height:     utils.Ptr(170.0),
			GoalWeight: utils.Ptr(60.0),
		}
		require.NoError(t, db.Create(&u).Error)

		ex := models.Exercise{Name: "Bench Press"}
		require.NoError(t, db.Create(&ex).Error)

		now := time.Now().Truncate(time.Second)
		recs := []models.WorkoutRecord{
			{UserID: u.ID, ExerciseID: ex.ID, BodyWeight: 61.0, TrainedOn: now},
			{UserID: u.ID, ExerciseID: ex.ID, BodyWeight: 62.0, TrainedOn: now},                   // 同日
			{UserID: u.ID, ExerciseID: ex.ID, BodyWeight: 63.0, TrainedOn: now.AddDate(0, 0, -1)}, // 別日
		}
		require.NoError(t, db.Create(&recs).Error)

		repo := repository.NewSummaryRepository(db)
		svc := service.NewSummaryService(repo)
		h := handler.NewSummaryHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/home/summary", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setUserID(c, u.ID)

		err := h.GetHomeSummary(c)
		if err != nil {
			e.HTTPErrorHandler(err, c)
		}

		require.Equal(t, http.StatusOK, rec.Code)

		var got SummaryResponseDTO
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))

		require.Equal(t, 2, got.TotalTrainingDays)
		require.NotNil(t, got.LatestWeight)
		require.InDelta(t, 62.0, *got.LatestWeight, 0.001)
		require.NotNil(t, got.LatestTrainedOn)
		expectedDate := now.Format("2006-01-02")
		require.Equal(t, expectedDate, got.LatestTrainedOn)
		require.NotNil(t, got.Height)
		require.InDelta(t, 170.0, *got.Height, 0.001)
		require.NotNil(t, got.GoalWeight)
		require.InDelta(t, 60.0, *got.GoalWeight, 0.001)
	})

	t.Run("【正常系】記録が存在しない場合は0件・最新体重nilで返ること", func(t *testing.T) {
		db := newSummaryIntegrationDB(t)

		u := models.User{
			Email:      "empty@example.com",
			Height:     utils.Ptr(180.0),
			GoalWeight: utils.Ptr(70.0),
		}
		require.NoError(t, db.Create(&u).Error)

		repo := repository.NewSummaryRepository(db)
		svc := service.NewSummaryService(repo)
		h := handler.NewSummaryHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/home/summary", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setUserID(c, u.ID)

		err := h.GetHomeSummary(c)
		if err != nil {
			e.HTTPErrorHandler(err, c)
		}

		require.Equal(t, http.StatusOK, rec.Code)

		var got service.HomeSummary
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))

		require.Equal(t, int64(0), got.TotalTrainingDays)
		require.Nil(t, got.LatestWeight)
		require.Nil(t, got.LatestTrainedOn)
		require.NotNil(t, got.Height)
		require.InDelta(t, 180.0, *got.Height, 0.001)
		require.NotNil(t, got.GoalWeight)
		require.InDelta(t, 70.0, *got.GoalWeight, 0.001)
	})

	t.Run("【異常系】内部エラー発生時は500(InternalError)が返ること", func(t *testing.T) {
		db := newSummaryIntegrationDB(t)

		u := models.User{Email: "boom@example.com"}
		require.NoError(t, db.Create(&u).Error)

		require.NoError(t, db.Migrator().DropTable(&models.WorkoutRecord{}))

		repo := repository.NewSummaryRepository(db)
		svc := service.NewSummaryService(repo)
		h := handler.NewSummaryHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/home/summary", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setUserID(c, u.ID)

		err := h.GetHomeSummary(c)
		if err != nil {
			e.HTTPErrorHandler(err, c)
		}

		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), `"code":"InternalError"`)
	})
}
