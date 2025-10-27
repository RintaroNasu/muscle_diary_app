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

	exRepo := repository.NewExerciseRepository(conn)
	exSvc := service.NewExerciseService(exRepo)
	exHandler := handler.NewExerciseHandler(exSvc)

	profileRepo := repository.NewProfileRepository(conn)
	profileSvc := service.NewProfileService(profileRepo)
	profileHandler := handler.NewProfileHandler(profileSvc)

	summaryRepo := repository.NewSummaryRepository(conn)
	summarySvc := service.NewSummaryService(summaryRepo)
	summaryHandler := handler.NewSummaryHandler(summarySvc)

	authRequired.GET("/exercises", exHandler.List)
	authRequired.POST("/training_records", workoutHandler.CreateWorkoutRecord)
	authRequired.GET("/training_records/date", workoutHandler.GetWorkoutRecordsByDate)
	authRequired.GET("/training_records/monthly_days", workoutHandler.GetMonthRecordDays)
	authRequired.PUT("/training_records/:id", workoutHandler.UpdateWorkoutRecord)
	authRequired.DELETE("/training_records/:id", workoutHandler.DeleteWorkoutRecord)
	authRequired.GET("/training_records/exercises/:exerciseId", workoutHandler.GetWorkoutRecordsByExercise)
	authRequired.GET("/profile", profileHandler.GetProfile)
	authRequired.PUT("/profile", profileHandler.UpdateProfile)
	authRequired.GET("/home/summary", summaryHandler.GetHomeSummary)
}
