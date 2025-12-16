package service

import "errors"

// 共通（Userドメインで再利用可能）
var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Workoutドメインで利用可能
var (
	ErrNoSets                 = errors.New("no sets")
	ErrInvalidSetValue        = errors.New("invalid set value")
	ErrExerciseNotFound       = errors.New("exercise not found")
	ErrRecordNotFound         = errors.New("record not found")
	ErrForbiddenPrivateRecord = errors.New("forbidden private record")
)
