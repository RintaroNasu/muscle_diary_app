package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RintaroNasu/muscle_diary_app/internal/handler"
	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAuthIntegration_SignUpAndLogin(t *testing.T) {
	t.Setenv("JWT_SECRET", "supersecret")
	e := echo.New()
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)
	svc := service.NewAuthService(repo)
	h := handler.NewAuthHandler(svc)

	// SignUp
	reqBody := `{"email":"test@test.com","password":"abcdef"}`
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := h.SignUp(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)

	// Login
	req2 := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)
	err = h.Login(c2)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec2.Code)
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}))
	return db
}
