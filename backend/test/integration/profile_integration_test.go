package integration

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func newProfileIntegrationDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}))
	return db
}

func setUserID(c echo.Context, id uint) {
	c.Set("user_id", id)
}

func TestProfileIntegration_GetAndUpdate(t *testing.T) {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	e.HTTPErrorHandler = httpx.HTTPErrorHandler(logger)

	t.Run("【正常系】プロフィール情報を取得できること", func(t *testing.T) {
		db := newProfileIntegrationDB(t)
		user := models.User{
			Email:      "user@example.com",
			Height:     utils.Ptr(170.0),
			GoalWeight: utils.Ptr(60.0),
		}
		require.NoError(t, db.Create(&user).Error)

		repo := repository.NewProfileRepository(db)
		svc := service.NewProfileService(repo)
		h := handler.NewProfileHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setUserID(c, user.ID)

		err := h.GetProfile(c)
		if err != nil {
			e.HTTPErrorHandler(err, c)
		}

		require.Equal(t, http.StatusOK, rec.Code)

		var res handler.ProfileResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		require.Equal(t, "user@example.com", res.Email)
		require.InDelta(t, 170.0, *res.HeightCM, 0.01)
		require.InDelta(t, 60.0, *res.GoalWeightKG, 0.01)
	})

	t.Run("【正常系】プロフィールを更新して最新の情報が返ること", func(t *testing.T) {
		db := newProfileIntegrationDB(t)
		user := models.User{
			Email:      "user2@example.com",
			Height:     utils.Ptr(160.0),
			GoalWeight: utils.Ptr(55.0),
		}
		require.NoError(t, db.Create(&user).Error)

		repo := repository.NewProfileRepository(db)
		svc := service.NewProfileService(repo)
		h := handler.NewProfileHandler(svc)

		body := `{"height_cm":175.0,"goal_weight_kg":65.0}`
		req := httptest.NewRequest(http.MethodPut, "/profile", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		setUserID(c, user.ID)

		err := h.UpdateProfile(c)
		if err != nil {
			e.HTTPErrorHandler(err, c)
		}

		require.Equal(t, http.StatusOK, rec.Code)

		var res handler.ProfileResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
		require.Equal(t, "user2@example.com", res.Email)
		require.InDelta(t, 175.0, *res.HeightCM, 0.01)
		require.InDelta(t, 65.0, *res.GoalWeightKG, 0.01)

		var updated models.User
		require.NoError(t, db.First(&updated, user.ID).Error)
		require.InDelta(t, 175.0, *updated.Height, 0.01)
		require.InDelta(t, 65.0, *updated.GoalWeight, 0.01)
	})

}
