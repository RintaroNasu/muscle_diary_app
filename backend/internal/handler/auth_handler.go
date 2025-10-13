package handler

import (
	"errors"
	"log/slog"
	"strings"

	"net/http"

	"github.com/RintaroNasu/muscle_diary_app/internal/httpx"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler interface {
	SignUp(c echo.Context) error
	Login(c echo.Context) error
}

type authHandler struct {
	svc service.AuthService
}

type authReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(svc service.AuthService) AuthHandler {
	return &authHandler{svc: svc}
}

func (h *authHandler) SignUp(c echo.Context) error {
	var req authReq
	ctx := c.Request().Context()

	if err := c.Bind(&req); err != nil {
		return httpx.BadRequest("InvalidBody", "リクエストの形式が不正です", err)
	}

	if req.Email == "" || req.Password == "" {
		return httpx.BadRequest("ValidationError", "email と password は必須です", nil)
	}

	if !strings.Contains(req.Email, "@") {
		return httpx.BadRequest("ValidationError", "email の形式が不正です", nil)
	}

	if len(req.Password) < 6 {
		return httpx.BadRequest("ValidationError", "password は6文字以上にしてください", nil)
	}

	u, token, err := h.svc.Signup(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			return httpx.Conflict("UserAlreadyExists", "すでに登録されています", err)
		}
		return httpx.Internal("ユーザー登録に失敗しました", err)
	}

	slog.InfoContext(ctx, "auth_signup_success", "user_id", u.ID)

	return c.JSON(http.StatusCreated, map[string]any{
		"id":    u.ID,
		"email": u.Email,
		"token": token,
	})
}

func (h *authHandler) Login(c echo.Context) error {
	var req authReq
	ctx := c.Request().Context()

	if err := c.Bind(&req); err != nil {
		return httpx.BadRequest("InvalidBody", "リクエストの形式が不正です", err)
	}

	if req.Email == "" || req.Password == "" {
		return httpx.BadRequest("ValidationError", "email と password は必須です", nil)
	}

	if !strings.Contains(req.Email, "@") {
		return httpx.BadRequest("ValidationError", "email の形式が不正です", nil)
	}

	if len(req.Password) < 6 {
		return httpx.BadRequest("ValidationError", "password は6文字以上にしてください", nil)
	}

	u, token, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound),
			errors.Is(err, service.ErrInvalidCredentials):
			return httpx.Unauthorized("認証に失敗しました", err)
		default:
			return httpx.Internal("ログインに失敗しました", err)
		}
	}

	slog.InfoContext(ctx, "auth_login_success", "user_id", u.ID)

	return c.JSON(http.StatusOK, map[string]any{
		"id":    u.ID,
		"email": u.Email,
		"token": token,
	})
}
