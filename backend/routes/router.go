package routes

import (
	"net/http"

	"github.com/RintaroNasu/muscle_diary_app/internal/handler"
	"github.com/RintaroNasu/muscle_diary_app/internal/middleware"
	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Register(e *echo.Echo, conn *gorm.DB) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Echo!")
	})
	authRepo := repository.NewUserRepository(conn)
	authSvc := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authSvc)

	e.POST("/signup", authHandler.SignUp)
	e.POST("/login", authHandler.Login)

	authRequired := e.Group("", middleware.JWTMiddleware())

	workoutRepo := repository.NewWorkoutRepository(conn)
	workoutSvc := service.NewWorkoutService(workoutRepo)
	workoutHandler := handler.NewWorkoutHandler(workoutSvc)
	authRequired.POST("/training_records", workoutHandler.CreateWorkoutRecord)
}
