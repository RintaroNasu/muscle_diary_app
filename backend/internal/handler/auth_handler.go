package handler

import (
	"net/http"

	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler interface {
	SignUp(c echo.Context) error
}

type authHandler struct {
	svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) AuthHandler {
	return &authHandler{svc: svc}
}

type signupReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *authHandler) SignUp(c echo.Context) error {
	var req signupReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	u, token, err := h.svc.Signup(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]any{
		"id":    u.ID,
		"email": u.Email,
		"token": token,
	})
}
