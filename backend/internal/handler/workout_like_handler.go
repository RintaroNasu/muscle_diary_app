package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/RintaroNasu/muscle_diary_app/internal/httpx"
	"github.com/RintaroNasu/muscle_diary_app/internal/middleware"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
)

type WorkoutLikeHandler interface {
	Like(c echo.Context) error
	Unlike(c echo.Context) error
}

type workoutLikeHandler struct {
	svc service.WorkoutLikeService
}

func NewWorkoutLikeHandler(svc service.WorkoutLikeService) WorkoutLikeHandler {
	return &workoutLikeHandler{svc: svc}
}

type LikeResponse struct {
	RecordID uint `json:"record_id"`
	Liked    bool `json:"liked"`
}

func parseRecordID(c echo.Context) (uint, error) {
	raw := c.Param("recordId")
	id64, err := strconv.ParseUint(raw, 10, 32)
	if err != nil || id64 == 0 {
		return 0, httpx.BadRequest("InvalidRecordID", "record_id が不正です", err)
	}
	return uint(id64), nil
}

func (h *workoutLikeHandler) Like(c echo.Context) error {
	ctx := c.Request().Context()

	recordID, err := parseRecordID(c)
	if err != nil {
		return err
	}

	userID := middleware.GetUserID(c)

	if err := h.svc.Like(userID, recordID); err != nil {
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			return httpx.NotFound("RecordNotFound", "対象レコードが存在しません", err)
		case errors.Is(err, service.ErrForbiddenPrivateRecord):
			return httpx.Forbidden("非公開の投稿にはいいねできません", err)
		default:
			return httpx.Internal("システムエラーが発生しました", err)
		}
	}

	slog.InfoContext(ctx, "like_created",
		"record_id", recordID,
		"user_id", userID,
	)

	return c.JSON(http.StatusOK, LikeResponse{
		RecordID: recordID,
		Liked:    true,
	})
}

func (h *workoutLikeHandler) Unlike(c echo.Context) error {
	ctx := c.Request().Context()

	recordID, err := parseRecordID(c)
	if err != nil {
		return err
	}

	userID := middleware.GetUserID(c)

	if err := h.svc.Unlike(userID, recordID); err != nil {
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			return httpx.NotFound("RecordNotFound", "対象レコードが存在しません", err)
		case errors.Is(err, service.ErrForbiddenPrivateRecord):
			return httpx.Forbidden("非公開の投稿には解除できません", err)
		default:
			return httpx.Internal("システムエラーが発生しました", err)
		}
	}

	slog.InfoContext(ctx, "like_deleted",
		"record_id", recordID,
		"user_id", userID,
	)

	return c.JSON(http.StatusOK, LikeResponse{
		RecordID: recordID,
		Liked:    false,
	})
}
