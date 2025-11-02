package service

import (
	"errors"
	"testing"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepo struct {
	findByEmail func(email string) (*models.User, error)
	create      func(u *models.User) error
}

func (f *fakeUserRepo) FindByEmail(email string) (*models.User, error) { return f.findByEmail(email) }
func (f *fakeUserRepo) Create(u *models.User) error                    { return f.create(u) }

func TestAuthService_Signup(t *testing.T) {
	type input struct {
		email    string
		password string
		secret   string
	}

	tests := []struct {
		name        string
		in          input
		repo        fakeUserRepo
		wantErr     error
		errContains string
	}{
		{
			name: "【正常系】ユーザーを新規登録できること",
			in:   input{"test@test.com", "asdfasdf", "supersecret"},
			repo: fakeUserRepo{
				findByEmail: func(email string) (*models.User, error) {
					return nil, repository.ErrNotFound
				},
				create: func(u *models.User) error {
					return nil
				},
			},
			wantErr: nil,
		},
		{
			name: "【異常系】既に登録済みのユーザーの場合は ErrUserAlreadyExists を返すこと",
			in:   input{"test@test.com", "asdfasdf", "supersecret"},
			repo: fakeUserRepo{
				findByEmail: func(email string) (*models.User, error) {
					return &models.User{Email: email}, nil
				},
				create: func(u *models.User) error { return nil },
			},
			wantErr: ErrUserAlreadyExists,
		},
		{
			name: "【異常系】FindByEmail が失敗した場合は find user failed エラーを返すこと",
			in:   input{"test@test.com", "asdfasdf", "supersecret"},
			repo: fakeUserRepo{
				findByEmail: func(email string) (*models.User, error) { return nil, errors.New("db down") },
				create:      func(u *models.User) error { return nil },
			},
			errContains: "find user failed",
		},
		{
			name: "【異常系】Create時に一意制約エラーが発生した場合は ErrUserAlreadyExists を返すこと",
			in:   input{"test@test.com", "asdfasdf", "supersecret"},
			repo: fakeUserRepo{
				findByEmail: func(email string) (*models.User, error) { return nil, repository.ErrNotFound },
				create:      func(u *models.User) error { return repository.ErrUniqueViolation },
			},
			wantErr: ErrUserAlreadyExists,
		},
		{
			name: "【異常系】Create 時にその他のDBエラーが発生した場合は create user failed エラーになること",
			in:   input{"test@test.com", "asdfasdf", "supersecret"},
			repo: fakeUserRepo{
				findByEmail: func(email string) (*models.User, error) { return nil, repository.ErrNotFound },
				create:      func(u *models.User) error { return errors.New("insert failed") },
			},
			errContains: "create user failed",
		},
		{
			name: "【異常系】JWT_SECRET が未設定の場合は jwt generate failed エラーになること",
			in:   input{"test@test.com", "asdfasdf", ""},
			repo: fakeUserRepo{
				findByEmail: func(email string) (*models.User, error) { return nil, repository.ErrNotFound },
				create:      func(u *models.User) error { u.ID = 1; return nil },
			},
			errContains: "jwt generate failed",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("JWT_SECRET", tt.in.secret)

			svc := NewAuthService(&tt.repo)
			u, token, err := svc.Signup(tt.in.email, tt.in.password)

			switch {
			case tt.wantErr != nil:
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, u)
				require.Empty(t, token)
			case tt.errContains != "":
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errContains)
			default:
				require.NoError(t, err)
				require.NotNil(t, u)
				require.Equal(t, tt.in.email, u.Email)
				require.NotEmpty(t, token)

				require.NotEqual(t, tt.in.password, u.Password)
				require.NoError(t, bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(tt.in.password)))
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	hash := func(pw string) string {
		b, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
		return string(b)
	}

	type input struct {
		email    string
		password string
		secret   string
	}

	tests := []struct {
		name        string
		in          input
		repo        fakeUserRepo
		wantErr     error
		errContains string
	}{
		{
			name: "【正常系】正しい認証情報でログインできトークンが発行されること",
			in:   input{"test@test.com", "asdfasdf", "supersecret"},
			repo: fakeUserRepo{
				findByEmail: func(email string) (*models.User, error) {
					return &models.User{
						Email:    email,
						Password: hash("asdfasdf"),
					}, nil
				},
				create: func(u *models.User) error { return nil },
			},
			wantErr: nil,
		},
		{
			name: "【異常系】存在しないユーザーの場合は ErrUserNotFound を返すこと",
			in:   input{"test@test.com", "asdfasdf", "supersecret"},
			repo: fakeUserRepo{
				findByEmail: func(email string) (*models.User, error) { return nil, repository.ErrNotFound },
				create:      func(u *models.User) error { return nil },
			},
			wantErr: ErrUserNotFound,
		},
		{
			name: "【異常系】FindByEmail が失敗した場合は find user failed エラーを返すこと",
			in:   input{"test@test.com", "asdfasdf", "supersecret"},
			repo: fakeUserRepo{
				findByEmail: func(email string) (*models.User, error) { return nil, errors.New("db down") },
				create:      func(u *models.User) error { return nil },
			},
			errContains: "find user failed",
		},
		{
			name: "【異常系】パスワードが不一致の場合は ErrInvalidCredentials を返すこと",
			in:   input{"test@test.com", "wrong", "supersecret"},
			repo: fakeUserRepo{
				findByEmail: func(email string) (*models.User, error) {
					return &models.User{
						Email:    "u@x.com",
						Password: hash("correct"),
					}, nil
				},
				create: func(u *models.User) error { return nil },
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			name: "【異常系】JWT_SECRET が未設定の場合は jwt generate failed エラーになること",
			in:   input{"test@test.com", "asdfasdf", ""},
			repo: fakeUserRepo{
				findByEmail: func(email string) (*models.User, error) {
					return &models.User{
						Email:    "ok@x.com",
						Password: hash("asdfasdf"),
					}, nil
				},
				create: func(u *models.User) error { return nil },
			},
			errContains: "jwt generate failed",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("JWT_SECRET", tt.in.secret)
			svc := NewAuthService(&tt.repo)

			u, token, err := svc.Login(tt.in.email, tt.in.password)

			switch {
			case tt.wantErr != nil:
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, u)
				require.Empty(t, token)
			case tt.errContains != "":
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errContains)
			default:
				require.NoError(t, err)
				require.NotNil(t, u)
				require.Equal(t, tt.in.email, u.Email)
				require.NotEmpty(t, token)
			}
		})
	}
}
