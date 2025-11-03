package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RintaroNasu/muscle_diary_app/internal/httpx"
	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/RintaroNasu/muscle_diary_app/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type fakeProfileService struct {
	getFunc    func(userID uint) (*models.User, error)
	updateFunc func(userID uint, height *float64, goalWeight *float64) (*models.User, error)
}

func (f *fakeProfileService) GetProfile(userID uint) (*models.User, error) {
	return f.getFunc(userID)
}
func (f *fakeProfileService) UpdateProfile(userID uint, h *float64, g *float64) (*models.User, error) {
	return f.updateFunc(userID, h, g)
}

func newEchoWithErrHandler() *echo.Echo {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	e.HTTPErrorHandler = httpx.HTTPErrorHandler(logger)
	return e
}

func setUserID(c echo.Context, id uint) { c.Set("user_id", id) }

func TestNewProfileHandler(t *testing.T) {
	h := NewProfileHandler(&fakeProfileService{})
	require.NotNil(t, h)
	_, ok := h.(ProfileHandler)
	require.True(t, ok)
}

func TestProfileHandler_GetProfile(t *testing.T) {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	e.HTTPErrorHandler = httpx.HTTPErrorHandler(logger)

	tests := []struct {
		name        string
		mock        fakeProfileService
		wantStatus  int
		wantBodyHas string
	}{
		{
			name: "【正常系】プロフィールを取得できる",
			mock: fakeProfileService{
				getFunc: func(userID uint) (*models.User, error) {
					return &models.User{
						Email:      "u@test.com",
						Height:     utils.Ptr(170.0),
						GoalWeight: utils.Ptr(60.0),
					}, nil
				},
			},
			wantStatus:  http.StatusOK,
			wantBodyHas: `"email":"u@test.com"`,
		},
		{
			name: "【異常系】ユーザーが存在しない場合は404(UserNotFound)を返すこと",
			mock: fakeProfileService{
				getFunc: func(userID uint) (*models.User, error) {
					return nil, service.ErrUserNotFound
				},
			},
			wantStatus:  http.StatusNotFound,
			wantBodyHas: `"UserNotFound"`,
		},
		{
			name: "【異常系】内部エラーは500を返すこと",
			mock: fakeProfileService{
				getFunc: func(userID uint) (*models.User, error) {
					return nil, errors.New("db down")
				},
			},
			wantStatus:  http.StatusInternalServerError,
			wantBodyHas: `"InternalError"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/me", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			setUserID(c, 1)

			h := NewProfileHandler(&tt.mock)
			err := h.GetProfile(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, tt.wantStatus, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantBodyHas)
		})
	}
}

func TestProfileHandler_UpdateProfile(t *testing.T) {
	e := newEchoWithErrHandler()

	tests := []struct {
		name        string
		body        string
		mock        fakeProfileService
		wantStatus  int
		wantBodyHas string
	}{
		{
			name: "【正常系】プロフィールを更新できること",
			body: `{"height_cm":175.5,"goal_weight_kg":65.2}`,
			mock: fakeProfileService{
				updateFunc: func(userID uint, h *float64, g *float64) (*models.User, error) {
					return &models.User{
						Email:      "u@test.com",
						Height:     h,
						GoalWeight: g,
					}, nil
				},
			},
			wantStatus:  http.StatusOK,
			wantBodyHas: `"email":"u@test.com"`,
		},
		{
			name: "【異常系】リクエストボディが不正なら400(InvalidBody)を返すこと",
			body: `{"height_cm": 170.0`,
			mock: fakeProfileService{
				updateFunc: func(userID uint, h *float64, g *float64) (*models.User, error) {
					return nil, nil
				},
			},
			wantStatus:  http.StatusBadRequest,
			wantBodyHas: `"InvalidBody"`,
		},
		{
			name: "【異常系】ユーザーが存在しない場合は404(UserNotFound)を返すこと",
			body: `{"height_cm":170,"goal_weight_kg":60}`,
			mock: fakeProfileService{
				updateFunc: func(userID uint, h *float64, g *float64) (*models.User, error) {
					return nil, service.ErrUserNotFound
				},
			},
			wantStatus:  http.StatusNotFound,
			wantBodyHas: `"UserNotFound"`,
		},
		{
			name: "【異常系】内部エラーは500を返すこと",
			body: `{"height_cm":170,"goal_weight_kg":60}`,
			mock: fakeProfileService{
				updateFunc: func(userID uint, h *float64, g *float64) (*models.User, error) {
					return nil, errors.New("update failed")
				},
			},
			wantStatus:  http.StatusInternalServerError,
			wantBodyHas: `"InternalError"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPatch, "/me", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			setUserID(c, 1)

			h := NewProfileHandler(&tt.mock)
			err := h.UpdateProfile(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, tt.wantStatus, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantBodyHas)

			if tt.wantStatus == http.StatusOK {
				var res ProfileResponse
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
				require.NotNil(t, res.HeightCM)
				require.NotNil(t, res.GoalWeightKG)
			}
		})
	}
}
