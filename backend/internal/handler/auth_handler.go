package handler

import (
	"fmt"
	"net/http"

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
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("invalid body: %v", err),
		})
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

func (h *authHandler) Login(c echo.Context) error {
	var req authReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("invalid body: %v", err),
		})
	}

	u, token, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"id":    u.ID,
		"email": u.Email,
		"token": token,
	})
}
