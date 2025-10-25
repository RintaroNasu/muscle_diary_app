package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/RintaroNasu/muscle_diary_app/internal/httpx"
	"github.com/RintaroNasu/muscle_diary_app/internal/middleware"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
)

type ProfileHandler interface {
	GetProfile(c echo.Context) error
	UpdateProfile(c echo.Context) error
}

type profileHandler struct {
	svc service.ProfileService
}

func NewProfileHandler(svc service.ProfileService) ProfileHandler {
	return &profileHandler{svc: svc}
}

type ProfileResponse struct {
	HeightCM     *float64 `json:"height_cm"`
	GoalWeightKG *float64 `json:"goal_weight_kg"`
	Email        string   `json:"email"`
}

type UpdateProfileRequest struct {
	HeightCM     *float64 `json:"height_cm"`
	GoalWeightKG *float64 `json:"goal_weight_kg"`
}

func (h *profileHandler) GetProfile(c echo.Context) error {
	userID := middleware.GetUserID(c)
	user, err := h.svc.GetProfile(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return httpx.NotFound("UserNotFound", "ユーザーが存在しません", err)
		}
		return httpx.Internal("システムエラーが発生しました", err)
	}

	res := ProfileResponse{
		HeightCM:     user.Height,
		GoalWeightKG: user.GoalWeight,
		Email:        user.Email,
	}
	return c.JSON(http.StatusOK, res)
}

func (h *profileHandler) UpdateProfile(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)
	var req UpdateProfileRequest

	if err := c.Bind(&req); err != nil {
		return httpx.BadRequest("InvalidBody", "リクエストの形式が不正です", err)
	}

	user, err := h.svc.UpdateProfile(userID, req.HeightCM, req.GoalWeightKG)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return httpx.NotFound("UserNotFound", "ユーザーが存在しません", err)
		}
		return httpx.Internal("システムエラーが発生しました", err)
	}

	slog.InfoContext(ctx, "profile_updated",
		"user_id", user.ID,
		"height_cm", user.Height,
		"goal_weight_kg", user.GoalWeight,
	)

	res := ProfileResponse{
		HeightCM:     user.Height,
		GoalWeightKG: user.GoalWeight,
		Email:        user.Email,
	}
	return c.JSON(http.StatusOK, res)
}
