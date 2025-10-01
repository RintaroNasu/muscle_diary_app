package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Signup(email, password string) (*models.User, string, error)
	Login(email, password string) (*models.User, string, error)
	// Login(email, password string) (string, error)
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo: repo}
}

type authService struct {
	repo repository.UserRepository
}

func (s *authService) Signup(email, password string) (*models.User, string, error) {
	if _, err := s.repo.FindByEmail(email); err == nil {
		return nil, "", errors.New("email already in use")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("password hash failed: %w", err)
	}

	u := &models.User{Email: email, Password: string(hash)}
	if err := s.repo.Create(u); err != nil {
		return nil, "", fmt.Errorf("create user failed: %w", err)
	}

	token, err := generateJWT(u.ID)
	if err != nil {
		return nil, "", fmt.Errorf("jwt generate failed: %w", err)
	}
	return u, token, nil
}

func (s *authService) Login(email, password string) (*models.User, string, error) {
	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("find user failed: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, "", fmt.Errorf("password compare failed: %w", err)
	}

	token, err := generateJWT(u.ID)
	if err != nil {
		return nil, "", fmt.Errorf("jwt generate failed: %w", err)
	}
	return u, token, nil
}

func generateJWT(userID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("missing JWT_SECRET")
	}
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(2 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}
