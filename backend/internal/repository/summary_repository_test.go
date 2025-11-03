package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newSummaryTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}, &models.Exercise{}, &models.WorkoutRecord{}))
	return db
}

func TestSummaryRepository_CountTrainingDays(t *testing.T) {
	tests := []struct {
		name        string
		prepare     func(db *gorm.DB)
		userID      uint
		wantCount   int64
		expectError bool
	}{
		{
			name: "【正常系】同一ユーザーの異なる訓練日が正しくカウントされること",
			prepare: func(db *gorm.DB) {
				user := models.User{Email: "test@example.com"}
				require.NoError(t, db.Create(&user).Error)
				exercise := models.Exercise{Name: "ベンチプレス"}
				require.NoError(t, db.Create(&exercise).Error)
				records := []models.WorkoutRecord{
					{UserID: user.ID, ExerciseID: exercise.ID, TrainedOn: time.Now()},
					{UserID: user.ID, ExerciseID: exercise.ID, TrainedOn: time.Now()},
					{UserID: user.ID, ExerciseID: exercise.ID, TrainedOn: time.Now().AddDate(0, 0, -1)},
				}
				require.NoError(t, db.Create(&records).Error)
			},
			userID:      1,
			wantCount:   2,
			expectError: false,
		},
		{
			name:        "【正常系】訓練データが存在しない場合は0を返すこと",
			prepare:     func(db *gorm.DB) {},
			userID:      1,
			wantCount:   0,
			expectError: false,
		},
		{
			name: "【異常系】テーブルが存在しない場合はエラーを返すこと",
			prepare: func(db *gorm.DB) {
				require.NoError(t, db.Migrator().DropTable(&models.WorkoutRecord{}))
			},
			userID:      1,
			wantCount:   0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newSummaryTestDB(t)
			tt.prepare(db)

			repo := NewSummaryRepository(db)
			got, err := repo.CountTrainingDays(tt.userID)

			if tt.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantCount, got)
		})
	}
}

func TestSummaryRepository_GetLatestWeight(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		prepare     func(db *gorm.DB)
		userID      uint
		wantWeight  *float64
		expectNil   bool
		expectError bool
	}{
		{
			name: "【正常系】最新の日付の体重を取得できること",
			prepare: func(db *gorm.DB) {
				user := models.User{Email: "test@example.com"}
				require.NoError(t, db.Create(&user).Error)
				exercise := models.Exercise{Name: "ベンチプレス"}
				require.NoError(t, db.Create(&exercise).Error)
				records := []models.WorkoutRecord{
					{UserID: user.ID, ExerciseID: exercise.ID, BodyWeight: 60.5, TrainedOn: now.AddDate(0, 0, -3)},
					{UserID: user.ID, ExerciseID: exercise.ID, BodyWeight: 61.2, TrainedOn: now.AddDate(0, 0, -1)},
					{UserID: user.ID, ExerciseID: exercise.ID, BodyWeight: 62.0, TrainedOn: now},
				}
				require.NoError(t, db.Create(&records).Error)
			},
			userID:      1,
			wantWeight:  ptr(62.0),
			expectNil:   false,
			expectError: false,
		},
		{
			name:        "【正常系】該当ユーザーのレコードが存在しない場合はnilを返すこと",
			prepare:     func(db *gorm.DB) {},
			userID:      99,
			wantWeight:  nil,
			expectNil:   true,
			expectError: false,
		},
		{
			name: "【異常系】テーブルが存在しない場合はエラーを返すこと",
			prepare: func(db *gorm.DB) {
				require.NoError(t, db.Migrator().DropTable(&models.WorkoutRecord{}))
			},
			userID:      1,
			wantWeight:  nil,
			expectNil:   true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newSummaryTestDB(t)
			tt.prepare(db)

			repo := NewSummaryRepository(db)
			weight, date, err := repo.GetLatestWeight(tt.userID)

			if tt.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if tt.expectNil {
				require.Nil(t, weight)
				require.Nil(t, date)
				return
			}

			require.NotNil(t, weight)
			require.InDelta(t, *tt.wantWeight, *weight, 0.001)
			require.WithinDuration(t, now, *date, time.Second)
		})
	}
}

func TestSummaryRepository_GetProfileBasics(t *testing.T) {
	tests := []struct {
		name        string
		prepare     func(db *gorm.DB)
		userID      uint
		wantHeight  float64
		wantGoal    float64
		expectError bool
	}{
		{
			name: "【正常系】ユーザーの身長・目標体重を取得できること",
			prepare: func(db *gorm.DB) {
				user := models.User{
					Email:      "test@example.com",
					Height:     ptr(175.0),
					GoalWeight: ptr(65.0),
				}
				require.NoError(t, db.Create(&user).Error)
			},
			userID:      1,
			wantHeight:  175.0,
			wantGoal:    65.0,
			expectError: false,
		},
		{
			name:        "【異常系】存在しないユーザーIDを指定した場合、RecordNotFoundエラーとなること",
			prepare:     func(db *gorm.DB) {},
			userID:      99,
			expectError: true,
		},
		{
			name: "【異常系】テーブルが存在しない場合、エラーを返すこと",
			prepare: func(db *gorm.DB) {
				require.NoError(t, db.Migrator().DropTable(&models.User{}))
			},
			userID:      1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newSummaryTestDB(t)
			tt.prepare(db)

			repo := NewSummaryRepository(db)
			height, goal, err := repo.GetProfileBasics(tt.userID)

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, height)
				require.Nil(t, goal)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, height)
			require.NotNil(t, goal)
			require.InDelta(t, tt.wantHeight, *height, 0.001)
			require.InDelta(t, tt.wantGoal, *goal, 0.001)
		})
	}
}
