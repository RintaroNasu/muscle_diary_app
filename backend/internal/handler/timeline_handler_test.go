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

// TimelineService のモック
type mockTimelineService struct {
	GetTimelineFunc func(userID uint) ([]service.TimelineItem, error)
}

func (m *mockTimelineService) GetTimeline(userID uint) ([]service.TimelineItem, error) {
	return m.GetTimelineFunc(userID)
}

func TestTimelineHandler_GetTimeline(t *testing.T) {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	e.HTTPErrorHandler = httpx.HTTPErrorHandler(logger)

	now := time.Now()

	tests := []struct {
		name           string
		mockSvc        service.TimelineService
		wantStatusCode int
		wantBodyPart   string
	}{
		{
			name: "【正常系】タイムラインを取得できること",
			mockSvc: &mockTimelineService{
				GetTimelineFunc: func(userID uint) ([]service.TimelineItem, error) {
					return []service.TimelineItem{
						{
							RecordID:     1,
							UserID:       10,
							UserEmail:    "user@example.com",
							ExerciseName: "ベンチプレス",
							BodyWeight:   70.5,
							TrainedOn:    now,
							Comment:      "今日は自己ベスト！",
							LikedByMe:    false,
						},
					}, nil
				},
			},
			wantStatusCode: http.StatusOK,
			wantBodyPart:   `"exercise_name":"ベンチプレス"`,
		},
		{
			name: "【異常系】サービス層でエラーが返された場合、500が返ること",
			mockSvc: &mockTimelineService{
				GetTimelineFunc: func(userID uint) ([]service.TimelineItem, error) {
					return nil, errors.New("db error")
				},
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBodyPart:   `"code":"InternalError"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewTimelineHandler(tt.mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/timeline", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user_id", uint(1))

			err := h.GetTimeline(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, tt.wantStatusCode, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantBodyPart)

			if tt.wantStatusCode == http.StatusOK {
				var res []TimelineItemResponse
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
				require.Len(t, res, 1)
				require.Equal(t, uint(1), res[0].RecordID)
				require.Equal(t, uint(10), res[0].UserID)
				require.Equal(t, "user@example.com", res[0].UserEmail)
				require.Equal(t, "ベンチプレス", res[0].ExerciseName)
				require.Equal(t, "今日は自己ベスト！", res[0].Comment)

				require.NotEmpty(t, res[0].TrainedOn)
			}
		})
	}
}
