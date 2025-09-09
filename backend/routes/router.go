package routes

import (
	"net/http"

	"github.com/RintaroNasu/muscle_diary_app/internal/handler"
	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Register(e *echo.Echo, conn *gorm.DB) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Echo!")
	})
	repo := repository.NewUserRepository(conn)
	svc := service.NewAuthService(repo)
	h := handler.NewAuthHandler(svc)

	e.POST("/signup", h.SignUp)
}
