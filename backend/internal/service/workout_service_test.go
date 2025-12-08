package service

import (
	"errors"
	"testing"
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type fakeWorkoutRepo struct {
	createFn   func(rec *models.WorkoutRecord) error
	findDayFn  func(userID uint, day time.Time) ([]models.WorkoutRecord, error)
	findMonFn  func(userID uint, year, month int) ([]time.Time, error)
	findOneFn  func(id uint, userID uint) (*models.WorkoutRecord, error)
	updateFn   func(rec *models.WorkoutRecord) error
	deleteFn   func(id uint, userID uint) error
	findSetsFn func(userID uint, exerciseID uint) ([]repository.FlatWorkoutSet, error)
}

func (f *fakeWorkoutRepo) Create(rec *models.WorkoutRecord) error {
	return f.createFn(rec)
}
func (f *fakeWorkoutRepo) FindByUserAndDay(userID uint, day time.Time) ([]models.WorkoutRecord, error) {
	return f.findDayFn(userID, day)
}
func (f *fakeWorkoutRepo) FindRecordDaysInMonth(userID uint, year int, month int) ([]time.Time, error) {
	return f.findMonFn(userID, year, month)
}
func (f *fakeWorkoutRepo) FindByIDAndUserID(id uint, userID uint) (*models.WorkoutRecord, error) {
	return f.findOneFn(id, userID)
}
func (f *fakeWorkoutRepo) Update(rec *models.WorkoutRecord) error {
	return f.updateFn(rec)
}
func (f *fakeWorkoutRepo) Delete(id uint, userID uint) error {
	return f.deleteFn(id, userID)
}
func (f *fakeWorkoutRepo) FindSetsByUserAndExercise(userID uint, exerciseID uint) ([]repository.FlatWorkoutSet, error) {
	return f.findSetsFn(userID, exerciseID)
}

func TestNewWorkoutService(t *testing.T) {
	svc := NewWorkoutService(&fakeWorkoutRepo{
		createFn:   func(*models.WorkoutRecord) error { return nil },
		findDayFn:  func(uint, time.Time) ([]models.WorkoutRecord, error) { return nil, nil },
		findMonFn:  func(uint, int, int) ([]time.Time, error) { return nil, nil },
		findOneFn:  func(uint, uint) (*models.WorkoutRecord, error) { return &models.WorkoutRecord{}, nil },
		updateFn:   func(*models.WorkoutRecord) error { return nil },
		deleteFn:   func(uint, uint) error { return nil },
		findSetsFn: func(uint, uint) ([]repository.FlatWorkoutSet, error) { return nil, nil },
	})
	require.NotNil(t, svc)
	_, ok := svc.(WorkoutService)
	require.True(t, ok)
}

func TestWorkoutService_CreateWorkoutRecord(t *testing.T) {
	day := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		repo       fakeWorkoutRepo
		userID     uint
		bodyWeight float64
		exerciseID uint
		trainedOn  time.Time
		sets       []WorkoutSetData
		isPublic   bool
		comment    string
		wantErr    error
		wantErrSub string
		wantSetLen int
	}{
		{
			name: "【正常系】レコードとセットを作成できること",
			repo: fakeWorkoutRepo{
				createFn: func(rec *models.WorkoutRecord) error {
					rec.ID = 10
					return nil
				},
			},
			userID:     1,
			bodyWeight: 70,
			exerciseID: 2,
			trainedOn:  day,
			sets: []WorkoutSetData{
				{SetNo: 1, Reps: 10, ExerciseWeight: 50},
				{SetNo: 2, Reps: 8, ExerciseWeight: 55},
			},
			wantSetLen: 2,
		},
		{
			name:    "【異常系】セットが空の場合は ErrNoSets を返すこと",
			repo:    fakeWorkoutRepo{},
			userID:  1,
			sets:    []WorkoutSetData{},
			wantErr: ErrNoSets,
		},
		{
			name:   "【異常系】セット値が不正な場合は ErrInvalidSetValue を返すこと",
			repo:   fakeWorkoutRepo{},
			userID: 1,
			sets: []WorkoutSetData{
				{SetNo: 0, Reps: 10, ExerciseWeight: 50},
			},
			wantErr: ErrInvalidSetValue,
		},
		{
			name: "【異常系】repo が FK 違反を返した場合は ErrExerciseNotFound を返すこと",
			repo: fakeWorkoutRepo{
				createFn: func(*models.WorkoutRecord) error {
					return repository.ErrFKViolation
				},
			},
			userID:  1,
			sets:    []WorkoutSetData{{SetNo: 1, Reps: 10, ExerciseWeight: 40}},
			wantErr: ErrExerciseNotFound,
		},
		{
			name: "【異常系】その他の repo エラーは wrap されて返すこと",
			repo: fakeWorkoutRepo{
				createFn: func(*models.WorkoutRecord) error {
					return errors.New("db down")
				},
			},
			userID:     1,
			exerciseID: 2,
			trainedOn:  day,
			sets:       []WorkoutSetData{{SetNo: 1, Reps: 10, ExerciseWeight: 40}},
			wantErrSub: "create workout record failed: db down",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := NewWorkoutService(&tt.repo)
			got, err := svc.CreateWorkoutRecord(tt.userID, tt.bodyWeight, tt.exerciseID, tt.trainedOn, tt.sets, tt.isPublic, tt.comment)

			if tt.wantErr != nil || tt.wantErrSub != "" {
				require.Error(t, err)
				if tt.wantErr != nil {
					require.ErrorIs(t, err, tt.wantErr)
				}
				if tt.wantErrSub != "" {
					require.Contains(t, err.Error(), tt.wantErrSub)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, tt.userID, got.UserID)
			require.Equal(t, tt.exerciseID, got.ExerciseID)
			require.Equal(t, tt.bodyWeight, got.BodyWeight)
			require.Equal(t, tt.trainedOn, got.TrainedOn)
			require.Len(t, got.Sets, tt.wantSetLen)
		})
	}
}

func TestWorkoutService_GetDailyRecords(t *testing.T) {
	day := time.Date(2025, 10, 2, 0, 0, 0, 0, time.UTC)
	dummy := []models.WorkoutRecord{{Model: gorm.Model{ID: 1}}, {Model: gorm.Model{ID: 2}}}

	tests := []struct {
		name       string
		repo       fakeWorkoutRepo
		userID     uint
		day        time.Time
		wantLen    int
		wantErrSub string
	}{
		{
			name: "【正常系】日別レコードを取得できること",
			repo: fakeWorkoutRepo{
				findDayFn: func(uint, time.Time) ([]models.WorkoutRecord, error) { return dummy, nil },
			},
			userID:  1,
			day:     day,
			wantLen: 2,
		},
		{
			name: "【異常系】repo エラーが発生した場合は wrap されて返すこと",
			repo: fakeWorkoutRepo{
				findDayFn: func(uint, time.Time) ([]models.WorkoutRecord, error) { return nil, errors.New("boom") },
			},
			userID:     1,
			day:        day,
			wantErrSub: "fetch daily records failed: boom",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := NewWorkoutService(&tt.repo)
			got, err := svc.GetDailyRecords(tt.userID, tt.day)

			if tt.wantErrSub != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrSub)
				return
			}
			require.NoError(t, err)
			require.Len(t, got, tt.wantLen)
		})
	}
}

func TestWorkoutService_GetMonthRecordDays(t *testing.T) {
	days := []time.Time{
		time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 10, 3, 0, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name       string
		repo       fakeWorkoutRepo
		userID     uint
		year       int
		month      int
		wantLen    int
		wantErrSub string
	}{
		{
			name: "【正常系】月内の日付を取得できること",
			repo: fakeWorkoutRepo{
				findMonFn: func(uint, int, int) ([]time.Time, error) { return days, nil },
			},
			userID:  1,
			year:    2025,
			month:   10,
			wantLen: 2,
		},
		{
			name:       "【異常系】month が不正な場合は invalid month エラーを返すこと",
			repo:       fakeWorkoutRepo{},
			userID:     1,
			year:       2025,
			month:      13,
			wantErrSub: "invalid month: 13",
		},
		{
			name: "【異常系】repo エラーが発生した場合は wrap されて返すこと",
			repo: fakeWorkoutRepo{
				findMonFn: func(uint, int, int) ([]time.Time, error) { return nil, errors.New("db err") },
			},
			userID:     1,
			year:       2025,
			month:      10,
			wantErrSub: "fetch month record days failed: db err",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := NewWorkoutService(&tt.repo)
			got, err := svc.GetMonthRecordDays(tt.userID, tt.year, tt.month)

			if tt.wantErrSub != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrSub)
				return
			}
			require.NoError(t, err)
			require.Len(t, got, tt.wantLen)
		})
	}
}

func TestWorkoutService_UpdateWorkoutRecord(t *testing.T) {
	fixed := time.Date(2025, 10, 5, 0, 0, 0, 0, time.UTC)

	baseRecord := &models.WorkoutRecord{
		Model:      gorm.Model{ID: 100},
		UserID:     1,
		ExerciseID: 2,
		BodyWeight: 60,
		TrainedOn:  time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
		Sets:       []models.WorkoutSet{{SetNo: 1, Reps: 10, ExerciseWeight: 40}},
	}

	tests := []struct {
		name       string
		repo       fakeWorkoutRepo
		userID     uint
		recordID   uint
		bodyWeight float64
		exerciseID uint
		trainedOn  time.Time
		sets       []WorkoutSetData
		wantErr    error
		wantErrSub string
		wantSetLen int
	}{
		{
			name: "【正常系】レコードとセットを更新できること",
			repo: fakeWorkoutRepo{
				findOneFn: func(id uint, userID uint) (*models.WorkoutRecord, error) {
					c := *baseRecord
					return &c, nil
				},
				updateFn: func(rec *models.WorkoutRecord) error { return nil },
			},
			userID:     1,
			recordID:   100,
			bodyWeight: 66.6,
			exerciseID: 3,
			trainedOn:  fixed,
			sets:       []WorkoutSetData{{SetNo: 1, Reps: 8, ExerciseWeight: 50}},
			wantSetLen: 1,
		},
		{
			name:    "【異常系】セットが空の場合は ErrNoSets を返すこと",
			repo:    fakeWorkoutRepo{},
			sets:    []WorkoutSetData{},
			wantErr: ErrNoSets,
		},
		{
			name:    "【異常系】セット値が不正な場合は ErrInvalidSetValue を返すこと",
			repo:    fakeWorkoutRepo{},
			sets:    []WorkoutSetData{{SetNo: 0, Reps: 10, ExerciseWeight: 40}},
			wantErr: ErrInvalidSetValue,
		},
		{
			name: "【異常系】FindByIDAndUserID が ErrNotFound を返した場合は ErrRecordNotFound を返すこと",
			repo: fakeWorkoutRepo{
				findOneFn: func(uint, uint) (*models.WorkoutRecord, error) { return nil, repository.ErrNotFound },
			},
			userID:   1,
			recordID: 999,
			sets:     []WorkoutSetData{{SetNo: 1, Reps: 10, ExerciseWeight: 40}},
			wantErr:  ErrRecordNotFound,
		},
		{
			name: "【異常系】FindByIDAndUserID その他のエラーは wrap されて返すこと",
			repo: fakeWorkoutRepo{
				findOneFn: func(uint, uint) (*models.WorkoutRecord, error) { return nil, errors.New("db find") },
			},
			userID:     1,
			recordID:   100,
			sets:       []WorkoutSetData{{SetNo: 1, Reps: 10, ExerciseWeight: 40}},
			wantErrSub: "find workout record failed: db find",
		},
		{
			name: "【異常系】Update 時 FK 違反を返した場合は ErrExerciseNotFound を返すこと",
			repo: fakeWorkoutRepo{
				findOneFn: func(uint, uint) (*models.WorkoutRecord, error) {
					c := *baseRecord
					return &c, nil
				},
				updateFn: func(*models.WorkoutRecord) error { return repository.ErrFKViolation },
			},
			userID:     1,
			recordID:   100,
			exerciseID: 999,
			trainedOn:  fixed,
			sets:       []WorkoutSetData{{SetNo: 1, Reps: 10, ExerciseWeight: 40}},
			wantErr:    ErrExerciseNotFound,
		},
		{
			name: "【異常系】Update その他のエラーは wrap されて返すこと",
			repo: fakeWorkoutRepo{
				findOneFn: func(uint, uint) (*models.WorkoutRecord, error) {
					c := *baseRecord
					return &c, nil
				},
				updateFn: func(*models.WorkoutRecord) error { return errors.New("update boom") },
			},
			userID:     1,
			recordID:   100,
			sets:       []WorkoutSetData{{SetNo: 1, Reps: 10, ExerciseWeight: 40}},
			wantErrSub: "update workout record failed: update boom",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := NewWorkoutService(&tt.repo)
			got, err := svc.UpdateWorkoutRecord(tt.userID, tt.recordID, tt.bodyWeight, tt.exerciseID, tt.trainedOn, tt.sets)

			if tt.wantErr != nil || tt.wantErrSub != "" {
				require.Error(t, err)
				if tt.wantErr != nil {
					require.ErrorIs(t, err, tt.wantErr)
				}
				if tt.wantErrSub != "" {
					require.Contains(t, err.Error(), tt.wantErrSub)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, tt.bodyWeight, got.BodyWeight)
			require.Equal(t, tt.exerciseID, got.ExerciseID)
			require.Equal(t, tt.trainedOn, got.TrainedOn)
			require.Len(t, got.Sets, tt.wantSetLen)
		})
	}
}

func TestWorkoutService_DeleteWorkoutRecord(t *testing.T) {
	tests := []struct {
		name       string
		repo       fakeWorkoutRepo
		userID     uint
		recordID   uint
		wantErr    error
		wantErrSub string
	}{
		{
			name: "【正常系】レコードを削除できること",
			repo: fakeWorkoutRepo{
				findOneFn: func(uint, uint) (*models.WorkoutRecord, error) {
					return &models.WorkoutRecord{Model: gorm.Model{ID: 10}, UserID: 1}, nil
				},
				deleteFn: func(uint, uint) error { return nil },
			},
			userID:   1,
			recordID: 10,
		},
		{
			name: "【異常系】Find が ErrNotFound を返した場合は ErrRecordNotFound を返すこと",
			repo: fakeWorkoutRepo{
				findOneFn: func(uint, uint) (*models.WorkoutRecord, error) { return nil, repository.ErrNotFound },
			},
			userID:   1,
			recordID: 999,
			wantErr:  ErrRecordNotFound,
		},
		{
			name: "【異常系】repo エラーが発生した場合は wrap されて返すこと",
			repo: fakeWorkoutRepo{
				findOneFn: func(uint, uint) (*models.WorkoutRecord, error) { return nil, errors.New("find boom") },
			},
			userID:     1,
			recordID:   10,
			wantErrSub: "find workout record failed: find boom",
		},
		{
			name: "【異常系】repo エラーが発生した場合は wrap されて返すこと",
			repo: fakeWorkoutRepo{
				findOneFn: func(uint, uint) (*models.WorkoutRecord, error) {
					return &models.WorkoutRecord{Model: gorm.Model{ID: 10}, UserID: 1}, nil
				},
				deleteFn: func(uint, uint) error { return errors.New("delete boom") },
			},
			userID:     1,
			recordID:   10,
			wantErrSub: "delete workout record failed: delete boom",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := NewWorkoutService(&tt.repo)
			err := svc.DeleteWorkoutRecord(tt.userID, tt.recordID)

			if tt.wantErr != nil || tt.wantErrSub != "" {
				require.Error(t, err)
				if tt.wantErr != nil {
					require.ErrorIs(t, err, tt.wantErr)
				}
				if tt.wantErrSub != "" {
					require.Contains(t, err.Error(), tt.wantErrSub)
				}
				return
			}
			require.NoError(t, err)
		})
	}
}
