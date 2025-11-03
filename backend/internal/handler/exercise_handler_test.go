package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RintaroNasu/muscle_diary_app/internal/httpx"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

type fakeExerciseService struct {
	listFunc func(ctx context.Context) ([]service.ExerciseDTO, error)
}

func (f *fakeExerciseService) List(ctx context.Context) ([]service.ExerciseDTO, error) {
	return f.listFunc(ctx)
}

func TestExerciseHandler_List(t *testing.T) {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	e.HTTPErrorHandler = httpx.HTTPErrorHandler(logger)

	tests := []struct {
		name        string
		mockSvc     fakeExerciseService
		wantStatus  int
		wantCount   int
		wantContain []string
		wantErrKey  string
	}{
		{
			name: "【正常系】種目一覧を取得できること",
			mockSvc: fakeExerciseService{
				listFunc: func(ctx context.Context) ([]service.ExerciseDTO, error) {
					return []service.ExerciseDTO{
						{ID: 1, Name: "Bench Press"},
						{ID: 2, Name: "Squat"},
					}, nil
				},
			},
			wantStatus:  http.StatusOK,
			wantCount:   2,
			wantContain: []string{"Bench Press", "Squat"},
		},
		{
			name: "【正常系】0件の場合は空配列を返すこと",
			mockSvc: fakeExerciseService{
				listFunc: func(ctx context.Context) ([]service.ExerciseDTO, error) {
					return []service.ExerciseDTO{}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantCount:  0,
		},
		{
			name: "【異常系】サービスエラー時は InternalError を返すこと",
			mockSvc: fakeExerciseService{
				listFunc: func(ctx context.Context) ([]service.ExerciseDTO, error) {
					return nil, errors.New("db down")
				},
			},
			wantStatus: http.StatusInternalServerError,
			wantErrKey: "InternalError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/exercises", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := NewExerciseHandler(&tt.mockSvc)
			err := h.List(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantErrKey != "" {
				require.Contains(t, rec.Body.String(), tt.wantErrKey)
				return
			}

			var got []service.ExerciseDTO
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
			require.Len(t, got, tt.wantCount)

			for _, want := range tt.wantContain {
				found := false
				for _, x := range got {
					if x.Name == want {
						found = true
						break
					}
				}
				require.True(t, found, "want name %q not found in response", want)
			}
		})
	}
}
