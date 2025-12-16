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

func newWorkoutLikeTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(
		&models.User{},
		&models.Exercise{},
		&models.WorkoutRecord{},
		&models.WorkoutLike{},
	))

	return db
}

func seedUserExerciseRecord(t *testing.T, db *gorm.DB, isPublic bool) (user models.User, ex models.Exercise, rec models.WorkoutRecord) {
	t.Helper()

	user = models.User{Email: "user@example.com", Password: "hashed"}
	require.NoError(t, db.Create(&user).Error)

	ex = models.Exercise{Name: "ベンチプレス"}
	require.NoError(t, db.Create(&ex).Error)

	rec = models.WorkoutRecord{
		UserID:     user.ID,
		ExerciseID: ex.ID,
		BodyWeight: 70.0,
		TrainedOn:  time.Now(),
		IsPublic:   isPublic,
		Comment:    "test",
	}
	require.NoError(t, db.Create(&rec).Error)

	return user, ex, rec
}

func TestWorkoutLikeRepository_CreateLike(t *testing.T) {
	tests := []struct {
		name        string
		prepare     func(db *gorm.DB) (userID uint, recordID uint)
		expectError bool
	}{
		{
			name: "【正常系】初回いいねが作成できること",
			prepare: func(db *gorm.DB) (uint, uint) {
				user, _, rec := seedUserExerciseRecord(t, db, true)
				return user.ID, rec.ID
			},
			expectError: false,
		},
		{
			name: "【正常系】同じ投稿に2回いいねしても冪等で成功扱いになること",
			prepare: func(db *gorm.DB) (uint, uint) {
				user, _, rec := seedUserExerciseRecord(t, db, true)
				// 1回目
				require.NoError(t, db.Create(&models.WorkoutLike{
					UserID:   user.ID,
					RecordID: rec.ID,
				}).Error)
				return user.ID, rec.ID
			},
			expectError: false,
		},
		{
			name: "【異常系】workout_likes テーブルが存在しない場合はエラーになること",
			prepare: func(db *gorm.DB) (uint, uint) {
				user, _, rec := seedUserExerciseRecord(t, db, true)
				require.NoError(t, db.Migrator().DropTable(&models.WorkoutLike{}))
				return user.ID, rec.ID
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db := newWorkoutLikeTestDB(t)
			userID, recordID := tt.prepare(db)

			repo := NewWorkoutLikeRepository(db)
			err := repo.CreateLike(userID, recordID)

			if tt.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			var cnt int64
			require.NoError(t, db.Model(&models.WorkoutLike{}).
				Where("user_id = ? AND record_id = ?", userID, recordID).
				Count(&cnt).Error)
			require.Equal(t, int64(1), cnt)
		})
	}
}

func TestWorkoutLikeRepository_DeleteLike(t *testing.T) {
	tests := []struct {
		name        string
		prepare     func(db *gorm.DB) (userID uint, recordID uint)
		expectError bool
		afterCount  int64
	}{
		{
			name: "【正常系】存在するいいねを物理削除できること（Unscoped）",
			prepare: func(db *gorm.DB) (uint, uint) {
				user, _, rec := seedUserExerciseRecord(t, db, true)
				require.NoError(t, db.Create(&models.WorkoutLike{
					UserID:   user.ID,
					RecordID: rec.ID,
				}).Error)
				return user.ID, rec.ID
			},
			expectError: false,
			afterCount:  0,
		},
		{
			name: "【正常系】存在しないいいねを削除しても冪等で成功扱いになること",
			prepare: func(db *gorm.DB) (uint, uint) {
				user, _, rec := seedUserExerciseRecord(t, db, true)
				return user.ID, rec.ID
			},
			expectError: false,
			afterCount:  0,
		},
		{
			name: "【異常系】workout_likes テーブルが存在しない場合はエラーになること",
			prepare: func(db *gorm.DB) (uint, uint) {
				user, _, rec := seedUserExerciseRecord(t, db, true)
				require.NoError(t, db.Migrator().DropTable(&models.WorkoutLike{}))
				return user.ID, rec.ID
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db := newWorkoutLikeTestDB(t)
			userID, recordID := tt.prepare(db)

			repo := NewWorkoutLikeRepository(db)
			err := repo.DeleteLike(userID, recordID)

			if tt.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			var cnt int64
			require.NoError(t, db.Unscoped().Model(&models.WorkoutLike{}).
				Where("user_id = ? AND record_id = ?", userID, recordID).
				Count(&cnt).Error)
			require.Equal(t, tt.afterCount, cnt)
		})
	}
}

func TestWorkoutLikeRepository_IsRecordPublic(t *testing.T) {
	tests := []struct {
		name        string
		prepare     func(db *gorm.DB) (recordID uint)
		wantPublic  bool
		wantErr     error
		expectError bool
	}{
		{
			name: "【正常系】公開レコードは true を返すこと",
			prepare: func(db *gorm.DB) uint {
				_, _, rec := seedUserExerciseRecord(t, db, true)
				return rec.ID
			},
			wantPublic:  true,
			expectError: false,
		},
		{
			name: "【正常系】非公開レコードは false を返すこと",
			prepare: func(db *gorm.DB) uint {
				_, _, rec := seedUserExerciseRecord(t, db, false)
				return rec.ID
			},
			wantPublic:  false,
			expectError: false,
		},
		{
			name: "【正常系】存在しない recordID は ErrNotFound を返すこと",
			prepare: func(db *gorm.DB) uint {
				return 999999
			},
			wantPublic:  false,
			wantErr:     ErrNotFound,
			expectError: true,
		},
		{
			name: "【異常系】workout_records テーブルが存在しない場合はエラーになること",
			prepare: func(db *gorm.DB) uint {
				require.NoError(t, db.Migrator().DropTable(&models.WorkoutRecord{}))
				return 1
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db := newWorkoutLikeTestDB(t)
			recordID := tt.prepare(db)

			repo := NewWorkoutLikeRepository(db)
			pub, err := repo.IsRecordPublic(recordID)

			if tt.expectError {
				require.Error(t, err)
				if tt.wantErr != nil {
					require.ErrorIs(t, err, tt.wantErr)
				}
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantPublic, pub)
		})
	}
}

func TestWorkoutLikeRepository_IsLikedByMe(t *testing.T) {
	tests := []struct {
		name        string
		prepare     func(db *gorm.DB) (userID uint, recordID uint)
		wantLiked   bool
		expectError bool
	}{
		{
			name: "【正常系】いいねしていない場合は false を返すこと",
			prepare: func(db *gorm.DB) (uint, uint) {
				user, _, rec := seedUserExerciseRecord(t, db, true)
				return user.ID, rec.ID
			},
			wantLiked:   false,
			expectError: false,
		},
		{
			name: "【正常系】いいね済みの場合は true を返すこと",
			prepare: func(db *gorm.DB) (uint, uint) {
				user, _, rec := seedUserExerciseRecord(t, db, true)
				require.NoError(t, db.Create(&models.WorkoutLike{
					UserID:   user.ID,
					RecordID: rec.ID,
				}).Error)
				return user.ID, rec.ID
			},
			wantLiked:   true,
			expectError: false,
		},
		{
			name: "【異常系】workout_likes テーブルが存在しない場合はエラーになること",
			prepare: func(db *gorm.DB) (uint, uint) {
				user, _, rec := seedUserExerciseRecord(t, db, true)
				require.NoError(t, db.Migrator().DropTable(&models.WorkoutLike{}))
				return user.ID, rec.ID
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db := newWorkoutLikeTestDB(t)
			userID, recordID := tt.prepare(db)

			repo := NewWorkoutLikeRepository(db)
			liked, err := repo.IsLikedByMe(userID, recordID)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantLiked, liked)
		})
	}
}
