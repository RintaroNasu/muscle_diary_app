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
)

type fakeWorkoutLikeService struct {
	likeFunc   func(userID uint, recordID uint) error
	unlikeFunc func(userID uint, recordID uint) error

	likeCalled   int
	unlikeCalled int
}

func (f *fakeWorkoutLikeService) Like(userID uint, recordID uint) error {
	f.likeCalled++
	if f.likeFunc == nil {
		return nil
	}
	return f.likeFunc(userID, recordID)
}

func (f *fakeWorkoutLikeService) Unlike(userID uint, recordID uint) error {
	f.unlikeCalled++
	if f.unlikeFunc == nil {
		return nil
	}
	return f.unlikeFunc(userID, recordID)
}

func TestNewWorkoutLikeHandler(t *testing.T) {
	h := NewWorkoutLikeHandler(&fakeWorkoutLikeService{})
	require.NotNil(t, h)
	_, ok := h.(WorkoutLikeHandler)
	require.True(t, ok)
}

func TestWorkoutLikeHandler_Like(t *testing.T) {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	e.HTTPErrorHandler = httpx.HTTPErrorHandler(logger)

	tests := []struct {
		name        string
		recordParam string
		mock        fakeWorkoutLikeService
		wantStatus  int
		wantBodyHas string
		wantLiked   *bool
		wantCalled  int
	}{
		{
			name:        "【正常系】いいねできること",
			recordParam: "10",
			mock: fakeWorkoutLikeService{
				likeFunc: func(userID uint, recordID uint) error {
					require.Equal(t, uint(1), userID)
					require.Equal(t, uint(10), recordID)
					return nil
				},
			},
			wantStatus:  http.StatusOK,
			wantBodyHas: `"record_id":10`,
			wantLiked:   func() *bool { b := true; return &b }(),
			wantCalled:  1,
		},
		{
			name:        "【異常系】record_id が数値でない場合は400(InvalidRecordID)",
			recordParam: "abc",
			mock: fakeWorkoutLikeService{
				likeFunc: func(userID uint, recordID uint) error { return nil },
			},
			wantStatus:  http.StatusBadRequest,
			wantBodyHas: `"InvalidRecordID"`,
			wantCalled:  0,
		},
		{
			name:        "【異常系】record_id が 0 の場合は400(InvalidRecordID)",
			recordParam: "0",
			mock: fakeWorkoutLikeService{
				likeFunc: func(userID uint, recordID uint) error { return nil },
			},
			wantStatus:  http.StatusBadRequest,
			wantBodyHas: `"InvalidRecordID"`,
			wantCalled:  0,
		},
		{
			name:        "【異常系】存在しないレコードは404(RecordNotFound)",
			recordParam: "999999",
			mock: fakeWorkoutLikeService{
				likeFunc: func(userID uint, recordID uint) error {
					return service.ErrRecordNotFound
				},
			},
			wantStatus:  http.StatusNotFound,
			wantBodyHas: `"RecordNotFound"`,
			wantCalled:  1,
		},
		{
			name:        "【異常系】非公開レコードは403",
			recordParam: "10",
			mock: fakeWorkoutLikeService{
				likeFunc: func(userID uint, recordID uint) error {
					return service.ErrForbiddenPrivateRecord
				},
			},
			wantStatus:  http.StatusForbidden,
			wantBodyHas: "非公開",
			wantCalled:  1,
		},
		{
			name:        "【異常系】想定外エラーは500(InternalError)",
			recordParam: "10",
			mock: fakeWorkoutLikeService{
				likeFunc: func(userID uint, recordID uint) error {
					return errors.New("db down")
				},
			},
			wantStatus:  http.StatusInternalServerError,
			wantBodyHas: `"InternalError"`,
			wantCalled:  1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/timeline/:recordId/like", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetParamNames("recordId")
			c.SetParamValues(tt.recordParam)

			c.Set("user_id", uint(1))

			h := NewWorkoutLikeHandler(&tt.mock)

			err := h.Like(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, tt.wantStatus, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantBodyHas)
			require.Equal(t, tt.wantCalled, tt.mock.likeCalled)

			if tt.wantStatus == http.StatusOK && tt.wantLiked != nil {
				var res LikeResponse
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
				require.Equal(t, uint(10), res.RecordID)
				require.Equal(t, *tt.wantLiked, res.Liked)
			}
		})
	}
}

func TestWorkoutLikeHandler_Unlike(t *testing.T) {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	e.HTTPErrorHandler = httpx.HTTPErrorHandler(logger)

	tests := []struct {
		name        string
		recordParam string
		mock        fakeWorkoutLikeService
		wantStatus  int
		wantBodyHas string
		wantLiked   *bool
		wantCalled  int
	}{
		{
			name:        "【正常系】いいね解除できること",
			recordParam: "10",
			mock: fakeWorkoutLikeService{
				unlikeFunc: func(userID uint, recordID uint) error {
					require.Equal(t, uint(1), userID)
					require.Equal(t, uint(10), recordID)
					return nil
				},
			},
			wantStatus:  http.StatusOK,
			wantBodyHas: `"record_id":10`,
			wantLiked:   func() *bool { b := false; return &b }(),
			wantCalled:  1,
		},
		{
			name:        "【異常系】record_id が不正なら400(InvalidRecordID)",
			recordParam: "abc",
			mock: fakeWorkoutLikeService{
				unlikeFunc: func(userID uint, recordID uint) error { return nil },
			},
			wantStatus:  http.StatusBadRequest,
			wantBodyHas: `"InvalidRecordID"`,
			wantCalled:  0,
		},
		{
			name:        "【異常系】存在しないレコードは404(RecordNotFound)",
			recordParam: "999999",
			mock: fakeWorkoutLikeService{
				unlikeFunc: func(userID uint, recordID uint) error {
					return service.ErrRecordNotFound
				},
			},
			wantStatus:  http.StatusNotFound,
			wantBodyHas: `"RecordNotFound"`,
			wantCalled:  1,
		},
		{
			name:        "【異常系】非公開レコードは403",
			recordParam: "10",
			mock: fakeWorkoutLikeService{
				unlikeFunc: func(userID uint, recordID uint) error {
					return service.ErrForbiddenPrivateRecord
				},
			},
			wantStatus:  http.StatusForbidden,
			wantBodyHas: "非公開",
			wantCalled:  1,
		},
		{
			name:        "【異常系】想定外エラーは500(InternalError)",
			recordParam: "10",
			mock: fakeWorkoutLikeService{
				unlikeFunc: func(userID uint, recordID uint) error {
					return errors.New("db down")
				},
			},
			wantStatus:  http.StatusInternalServerError,
			wantBodyHas: `"InternalError"`,
			wantCalled:  1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/timeline/:recordId/like", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetParamNames("recordId")
			c.SetParamValues(tt.recordParam)

			c.Set("user_id", uint(1))

			h := NewWorkoutLikeHandler(&tt.mock)

			err := h.Unlike(c)
			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, tt.wantStatus, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantBodyHas)
			require.Equal(t, tt.wantCalled, tt.mock.unlikeCalled)

			if tt.wantStatus == http.StatusOK && tt.wantLiked != nil {
				var res LikeResponse
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
				require.Equal(t, uint(10), res.RecordID)
				require.Equal(t, *tt.wantLiked, res.Liked)
			}
		})
	}
}
