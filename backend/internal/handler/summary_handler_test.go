package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/httpx"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type mockSummaryService struct {
	GetHomeSummaryFunc func(userID uint) (*service.HomeSummary, error)
}

func (m *mockSummaryService) GetHomeSummary(userID uint) (*service.HomeSummary, error) {
	return m.GetHomeSummaryFunc(userID)
}

func TestSummaryHandler_GetHomeSummary(t *testing.T) {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	e.HTTPErrorHandler = httpx.HTTPErrorHandler(logger)

	now := time.Now()

	tests := []struct {
		name           string
		mockSvc        service.SummaryService
		wantStatusCode int
		wantBodyPart   string
	}{
		{
			name: "【正常系】サマリーを取得できること",
			mockSvc: &mockSummaryService{
				GetHomeSummaryFunc: func(userID uint) (*service.HomeSummary, error) {
					return &service.HomeSummary{
						TotalTrainingDays: 10,
						LatestWeight:      ptr(65.0),
						LatestTrainedOn:   &now,
						GoalWeight:        ptr(60.0),
						Height:            ptr(170.0),
					}, nil
				},
			},
			wantStatusCode: http.StatusOK,
			wantBodyPart:   `"total_training_days":10`,
		},
		{
			name: "【異常系】サービス層でエラーが返された場合、500が返ること",
			mockSvc: &mockSummaryService{
				GetHomeSummaryFunc: func(userID uint) (*service.HomeSummary, error) {
					return nil, errors.New("test error")
				},
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBodyPart:   `"code":"InternalError"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewSummaryHandler(tt.mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/home/summary", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.Set("user_id", uint(1))

			err := h.GetHomeSummary(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, tt.wantStatusCode, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantBodyPart)

			if tt.wantStatusCode == http.StatusOK {
				var res service.HomeSummary
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
				require.Equal(t, int64(10), res.TotalTrainingDays)
				require.InDelta(t, 65.0, *res.LatestWeight, 0.01)
			}
		})
	}
}
