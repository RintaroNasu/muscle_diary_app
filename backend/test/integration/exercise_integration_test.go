package integration

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
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

func newIntegrationDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.Exercise{}))
	return db
}

func TestExerciseIntegration_List(t *testing.T) {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(nil, &slog.HandlerOptions{}))
	e.HTTPErrorHandler = httpx.HTTPErrorHandler(logger)

	t.Run("【正常系】DBのデータをAPI経由で取得できること", func(t *testing.T) {
		db := newIntegrationDB(t)
		seeds := []models.Exercise{
			{Name: "Bench Press"},
			{Name: "Squat"},
			{Name: "Deadlift"},
		}
		require.NoError(t, db.Create(&seeds).Error)

		repo := repository.NewExerciseRepository(db)
		svc := service.NewExerciseService(repo)
		h := handler.NewExerciseHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/exercises", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.List(c)
		if err != nil {
			e.HTTPErrorHandler(err, c)
		}

		require.Equal(t, http.StatusOK, rec.Code)

		var got []service.ExerciseDTO
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
		require.Len(t, got, 3)
		require.Equal(t, "Bench Press", got[0].Name)
		require.Equal(t, "Squat", got[1].Name)
		require.Equal(t, "Deadlift", got[2].Name)
	})
}
