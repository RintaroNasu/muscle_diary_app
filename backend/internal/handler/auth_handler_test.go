package handler

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RintaroNasu/muscle_diary_app/internal/httpx"
	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/RintaroNasu/muscle_diary_app/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type fakeAuthService struct {
	signupFunc func(email, password string) (*models.User, string, error)
	loginFunc  func(email, password string) (*models.User, string, error)
}

func (f *fakeAuthService) Signup(email, password string) (*models.User, string, error) {
	return f.signupFunc(email, password)
}
func (f *fakeAuthService) Login(email, password string) (*models.User, string, error) {
	return f.loginFunc(email, password)
}

func TestAuthHandler_SignUp(t *testing.T) {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	e.HTTPErrorHandler = httpx.HTTPErrorHandler(logger)

	tests := []struct {
		name       string
		body       string
		mockSvc    fakeAuthService
		wantStatus int
		wantBody   string
	}{
		{
			name: "【正常系】ユーザーを新規登録できること",
			body: `{"email":"test@test.com","password":"asdfasdf"}`,
			mockSvc: fakeAuthService{
				signupFunc: func(email, password string) (*models.User, string, error) {
					return &models.User{Email: email}, "token", nil
				},
			},
			wantStatus: http.StatusCreated,
			wantBody:   `"token":"token"`,
		},
		{
			name: "【異常系】既に登録済みのユーザーの場合は ErrUserAlreadyExists を返すこと",
			body: `{"email":"dup@test.com","password":"asdfasdf"}`,
			mockSvc: fakeAuthService{
				signupFunc: func(email, password string) (*models.User, string, error) {
					return nil, "", service.ErrUserAlreadyExists
				},
			},
			wantStatus: http.StatusConflict,
			wantBody:   "UserAlreadyExists",
		},
		{
			name: "【異常系】不正なリクエストの場合は InvalidBody エラーを返すこと",
			body: `{"email": "broken"`,
			mockSvc: fakeAuthService{
				signupFunc: func(email, password string) (*models.User, string, error) {
					return nil, "", nil
				},
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "InvalidBody",
		},
		{
			name:       "【異常系】パスワードが短い場合は password too short エラーを返すこと",
			body:       `{"email":"test@test.com","password":"123"}`,
			mockSvc:    fakeAuthService{},
			wantStatus: http.StatusBadRequest,
			wantBody:   `"password は6文字以上にしてください"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := &authHandler{svc: &tt.mockSvc}
			err := h.SignUp(c)

			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, tt.wantStatus, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantBody)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	e := echo.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	e.HTTPErrorHandler = httpx.HTTPErrorHandler(logger)

	tests := []struct {
		name       string
		body       string
		mockSvc    fakeAuthService
		wantStatus int
		wantBody   string
	}{
		{
			name: "【正常系】正しい認証情報でログインできること",
			body: `{"email":"test@test.com","password":"asdfasdf"}`,
			mockSvc: fakeAuthService{
				loginFunc: func(email, password string) (*models.User, string, error) {
					return &models.User{Email: email}, "token", nil
				},
			},
			wantStatus: http.StatusOK,
			wantBody:   `"token":"token"`,
		},
		{
			name: "【異常系】存在しないユーザーでログインした場合は UserNotFound エラーを返すこと",
			body: `{"email":"none@test.com","password":"asdfasdf"}`,
			mockSvc: fakeAuthService{
				loginFunc: func(email, password string) (*models.User, string, error) {
					return nil, "", service.ErrUserNotFound
				},
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   `"認証に失敗しました"`,
		},
		{
			name: "【異常系】パスワードが不一致の場合は InvalidCredentials エラーを返すこと",
			body: `{"email":"test@test.com","password":"wrongpw"}`,
			mockSvc: fakeAuthService{
				loginFunc: func(email, password string) (*models.User, string, error) {
					return nil, "", service.ErrInvalidCredentials
				},
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   `"認証に失敗しました"`,
		},
		{
			name: "【異常系】サービス側の内部エラーが発生した場合は InternalError エラーを返すこと",
			body: `{"email":"test@test.com","password":"asdfasdf"}`,
			mockSvc: fakeAuthService{
				loginFunc: func(email, password string) (*models.User, string, error) {
					return nil, "", errors.New("db down")
				},
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "InternalError",
		},
		{
			name:       "【異常系】不正なemail形式の場合は ValidationError エラーを返すこと",
			body:       `{"email":"invalid","password":"asdfasdf"}`,
			mockSvc:    fakeAuthService{},
			wantStatus: http.StatusBadRequest,
			wantBody:   `"email の形式が不正です"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := &authHandler{svc: &tt.mockSvc}
			err := h.Login(c)

			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			require.Equal(t, tt.wantStatus, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantBody)
		})
	}
}
