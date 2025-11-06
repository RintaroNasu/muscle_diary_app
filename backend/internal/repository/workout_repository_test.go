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

func newWorkoutTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(
		&models.User{},
		&models.Exercise{},
		&models.WorkoutRecord{},
		&models.WorkoutSet{},
	))
	return db
}

func TestWorkoutRepository_Create(t *testing.T) {
	tests := []struct {
		name        string
		prepare     func(db *gorm.DB)
		record      models.WorkoutRecord
		expectError bool
	}{
		{
			name: "【正常系】レコードを作成できること",
			prepare: func(db *gorm.DB) {
				require.NoError(t, db.Create(&models.User{Email: "test@example.com"}).Error)
				require.NoError(t, db.Create(&models.Exercise{Name: "ベンチプレス"}).Error)
			},
			record: models.WorkoutRecord{
				BodyWeight: 70,
				TrainedOn:  time.Now(),
			},
			expectError: false,
		},
		{
			name: "【異常系】存在しないUserIDを指定した場合 foreign_key エラーを返すこと",
			prepare: func(db *gorm.DB) {
			},
			record: models.WorkoutRecord{
				UserID:     9999,
				ExerciseID: 9999,
				BodyWeight: 70,
				TrainedOn:  time.Now(),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newWorkoutTestDB(t)
			repo := NewWorkoutRepository(db)

			tt.prepare(db)
			rec := tt.record
			if rec.UserID == 0 {
				rec.UserID = 1
			}
			if rec.ExerciseID == 0 {
				rec.ExerciseID = 1
			}
			err := repo.Create(&rec)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotZero(t, rec.ID)
		})
	}
}
func TestWorkoutRepository_FindByUserAndDay(t *testing.T) {
	tests := []struct {
		name        string
		prepare     func(db *gorm.DB)
		userID      uint
		day         time.Time
		wantLen     int
		expectError bool
	}{
		{
			name: "【正常系】同一ユーザーの同一日付レコードのみ取得できること",
			prepare: func(db *gorm.DB) {
				u1 := models.User{Email: "u1@example.com"}
				u2 := models.User{Email: "u2@example.com"}
				require.NoError(t, db.Create(&u1).Error)
				require.NoError(t, db.Create(&u2).Error)

				ex := models.Exercise{Name: "ベンチプレス"}
				require.NoError(t, db.Create(&ex).Error)

				day := time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC)

				recs := []models.WorkoutRecord{
					{UserID: u1.ID, ExerciseID: ex.ID, TrainedOn: day, BodyWeight: 70},
					{UserID: u1.ID, ExerciseID: ex.ID, TrainedOn: day, BodyWeight: 71},
					{UserID: u1.ID, ExerciseID: ex.ID, TrainedOn: day.AddDate(0, 0, -1), BodyWeight: 69}, // 別日
					{UserID: u2.ID, ExerciseID: ex.ID, TrainedOn: day, BodyWeight: 60},                   // 別ユーザー
				}
				require.NoError(t, db.Create(&recs).Error)

				// Preload の確認用に set を1件だけ付与
				require.NoError(t, db.Create(&models.WorkoutSet{
					WorkoutRecordID: recs[0].ID, SetNo: 1, Reps: 10, ExerciseWeight: 40,
				}).Error)
			},
			userID:      1,
			day:         time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC),
			wantLen:     2,
			expectError: false,
		},
		{
			name: "【正常系】該当が無ければ空配列を返すこと",
			prepare: func(db *gorm.DB) {
				require.NoError(t, db.Create(&models.User{Email: "u@example.com"}).Error)
				require.NoError(t, db.Create(&models.Exercise{Name: "DL"}).Error)
			},
			userID:      1,
			day:         time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC),
			wantLen:     0,
			expectError: false,
		},
		{
			name: "【異常系】テーブルが存在しない場合はエラーを返すこと",
			prepare: func(db *gorm.DB) {
				require.NoError(t, db.Migrator().DropTable(&models.WorkoutRecord{}))
			},
			userID:      1,
			day:         time.Now(),
			wantLen:     0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db := newWorkoutTestDB(t)
			repo := NewWorkoutRepository(db)

			tt.prepare(db)

			got, err := repo.FindByUserAndDay(tt.userID, tt.day)
			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Len(t, got, tt.wantLen)
			if tt.wantLen > 0 {
				require.NotNil(t, got[0].Exercise)
				require.NotNil(t, got[0].Sets)
			}
		})
	}
}
func TestWorkoutRepository_FindByIDAndUserID(t *testing.T) {
	tests := []struct {
		name        string
		prepare     func(db *gorm.DB)
		id          uint
		userID      uint
		expectNil   bool
		expectError bool
	}{
		{
			name: "【正常系】ID と userID に一致するレコードを取得できること",
			prepare: func(db *gorm.DB) {
				user := models.User{Email: "test@example.com"}
				require.NoError(t, db.Create(&user).Error)

				exercise := models.Exercise{Name: "ベンチプレス"}
				require.NoError(t, db.Create(&exercise).Error)

				record := models.WorkoutRecord{
					UserID:     user.ID,
					ExerciseID: exercise.ID,
					BodyWeight: 70.0,
					TrainedOn:  time.Now(),
				}
				require.NoError(t, db.Create(&record).Error)
			},
			expectNil:   false,
			expectError: false,
		},
		{
			name:        "【異常系】存在しないIDの場合 ErrNotFound を返すこと",
			prepare:     func(db *gorm.DB) {},
			id:          9999,
			userID:      1,
			expectNil:   true,
			expectError: true,
		},
		{
			name: "【異常系】テーブルが存在しない場合 エラーを返すこと",
			prepare: func(db *gorm.DB) {
				require.NoError(t, db.Migrator().DropTable(&models.WorkoutRecord{}))
			},
			id:          1,
			userID:      1,
			expectNil:   true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newWorkoutTestDB(t)
			repo := NewWorkoutRepository(db)

			tt.prepare(db)

			if tt.id == 0 && tt.userID == 0 && !tt.expectError {
				var r models.WorkoutRecord
				require.NoError(t, db.First(&r).Error)
				tt.id = r.ID
				tt.userID = r.UserID
			}

			got, err := repo.FindByIDAndUserID(tt.id, tt.userID)

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, tt.id, got.ID)
		})
	}
}

func TestWorkoutRepository_Update(t *testing.T) {
	tests := []struct {
		name        string
		prepare     func(db *gorm.DB) (recID uint, userID uint)
		mutate      func(db *gorm.DB, rec *models.WorkoutRecord)
		verify      func(t *testing.T, db *gorm.DB, recID uint)
		expectError bool
	}{
		{
			name: "【正常系】フィールド更新とセット差し替えができること",
			prepare: func(db *gorm.DB) (uint, uint) {
				u := models.User{Email: "upd@example.com"}
				ex1 := models.Exercise{Name: "ショルダー"}
				ex2 := models.Exercise{Name: "ラットプル"}
				require.NoError(t, db.Create(&u).Error)
				require.NoError(t, db.Create(&ex1).Error)
				require.NoError(t, db.Create(&ex2).Error)

				rec := models.WorkoutRecord{
					UserID:     u.ID,
					ExerciseID: ex1.ID,
					BodyWeight: 65,
					TrainedOn:  time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
					Sets: []models.WorkoutSet{
						{SetNo: 1, Reps: 12, ExerciseWeight: 20},
						{SetNo: 2, Reps: 10, ExerciseWeight: 22.5},
					},
				}
				require.NoError(t, db.Create(&rec).Error)
				return rec.ID, u.ID
			},
			mutate: func(db *gorm.DB, rec *models.WorkoutRecord) {
				var ex2 models.Exercise
				require.NoError(t, db.Where("name = ?", "ラットプル").First(&ex2).Error)

				rec.ExerciseID = ex2.ID
				rec.BodyWeight = 66.2
				rec.TrainedOn = time.Date(2025, 7, 2, 0, 0, 0, 0, time.UTC)
				rec.Sets = []models.WorkoutSet{
					{SetNo: 1, Reps: 8, ExerciseWeight: 30},
				}
			},
			verify: func(t *testing.T, db *gorm.DB, recID uint) {
				var got models.WorkoutRecord
				require.NoError(t, db.Preload("Sets").First(&got, recID).Error)
				require.InDelta(t, 66.2, got.BodyWeight, 0.001)
				require.Equal(t, time.Date(2025, 7, 2, 0, 0, 0, 0, time.UTC), got.TrainedOn)
				require.Len(t, got.Sets, 1)
				require.Equal(t, 1, got.Sets[0].SetNo)
				require.InDelta(t, 30.0, got.Sets[0].ExerciseWeight, 0.001)
			},
			expectError: false,
		},
		{
			name: "【異常系】存在しないexercise_idに更新するとエラーになること",
			prepare: func(db *gorm.DB) (uint, uint) {
				u := models.User{Email: "upd_fk@example.com"}
				ex := models.Exercise{Name: "アップライト"}
				require.NoError(t, db.Create(&u).Error)
				require.NoError(t, db.Create(&ex).Error)

				rec := models.WorkoutRecord{
					UserID:     u.ID,
					ExerciseID: ex.ID,
					TrainedOn:  time.Now(),
					Sets:       []models.WorkoutSet{{SetNo: 1, Reps: 10, ExerciseWeight: 20}},
				}
				require.NoError(t, db.Create(&rec).Error)
				return rec.ID, u.ID
			},
			mutate: func(db *gorm.DB, rec *models.WorkoutRecord) {
				rec.ExerciseID = 999999
			},
			verify:      func(t *testing.T, db *gorm.DB, recID uint) {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db := newWorkoutTestDB(t)
			repo := NewWorkoutRepository(db)

			recID, _ := tt.prepare(db)

			var rec models.WorkoutRecord
			require.NoError(t, db.Preload("Sets").First(&rec, recID).Error)

			tt.mutate(db, &rec)
			err := repo.Update(&rec)
			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			tt.verify(t, db, recID)
		})
	}
}

func TestWorkoutRepository_Delete(t *testing.T) {
	tests := []struct {
		name      string
		prepare   func(db *gorm.DB) (recID uint, ownerID uint, otherID uint)
		doDelete  func(repo WorkoutRepository, recID uint, ownerID uint, otherID uint) (deletedByOther bool, deletedByOwner bool, err error)
		verify    func(t *testing.T, db *gorm.DB, recID uint, deletedByOther, deletedByOwner bool)
		expectErr bool
	}{
		{
			name: "【正常系】本人のみ削除できること",
			prepare: func(db *gorm.DB) (uint, uint, uint) {
				u1 := models.User{Email: "owner@example.com"}
				u2 := models.User{Email: "other@example.com"}
				ex := models.Exercise{Name: "ローイング"}
				require.NoError(t, db.Create(&u1).Error)
				require.NoError(t, db.Create(&u2).Error)
				require.NoError(t, db.Create(&ex).Error)

				rec := models.WorkoutRecord{UserID: u1.ID, ExerciseID: ex.ID, TrainedOn: time.Now()}
				require.NoError(t, db.Create(&rec).Error)
				return rec.ID, u1.ID, u2.ID
			},
			doDelete: func(repo WorkoutRepository, recID uint, ownerID uint, otherID uint) (bool, bool, error) {
				err1 := repo.Delete(recID, otherID)
				err2 := repo.Delete(recID, ownerID)
				return err1 == nil, err2 == nil, err1
			},
			verify: func(t *testing.T, db *gorm.DB, recID uint, deletedByOther, deletedByOwner bool) {
				var cnt int64
				require.NoError(t, db.Model(&models.WorkoutRecord{}).Where("id = ?", recID).Count(&cnt).Error)
				require.True(t, deletedByOther)
				require.True(t, deletedByOwner)
				require.Equal(t, int64(0), cnt)
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db := newWorkoutTestDB(t)
			repo := NewWorkoutRepository(db)

			recID, ownerID, otherID := tt.prepare(db)
			deletedByOther, deletedByOwner, err := tt.doDelete(repo, recID, ownerID, otherID)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			tt.verify(t, db, recID, deletedByOther, deletedByOwner)
		})
	}
}

func TestWorkoutRepository_FindSetsByUserAndExercise(t *testing.T) {
	tests := []struct {
		name        string
		prepare     func(db *gorm.DB) (userID uint, ex1ID uint)
		userID      uint
		exerciseID  uint
		wantLen     int
		expectError bool
	}{
		{
			name: "【正常系】指定ユーザー×種目のセットが日付→セット番号順で取得できること",
			prepare: func(db *gorm.DB) (uint, uint) {
				u := models.User{Email: "sets@example.com"}
				ex1 := models.Exercise{Name: "スクワット"}
				ex2 := models.Exercise{Name: "デッドリフト"}
				require.NoError(t, db.Create(&u).Error)
				require.NoError(t, db.Create(&ex1).Error)
				require.NoError(t, db.Create(&ex2).Error)

				day1 := time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC)
				day2 := time.Date(2025, 9, 3, 0, 0, 0, 0, time.UTC)

				r1 := models.WorkoutRecord{UserID: u.ID, ExerciseID: ex1.ID, BodyWeight: 70, TrainedOn: day1}
				require.NoError(t, db.Create(&r1).Error)
				r2 := models.WorkoutRecord{UserID: u.ID, ExerciseID: ex1.ID, BodyWeight: 70.5, TrainedOn: day2}
				require.NoError(t, db.Create(&r2).Error)
				r3 := models.WorkoutRecord{UserID: u.ID, ExerciseID: ex2.ID, BodyWeight: 71, TrainedOn: day1}
				require.NoError(t, db.Create(&r3).Error)

				require.NoError(t, db.Create(&[]models.WorkoutSet{
					{WorkoutRecordID: r1.ID, SetNo: 2, Reps: 8, ExerciseWeight: 100},
					{WorkoutRecordID: r1.ID, SetNo: 1, Reps: 10, ExerciseWeight: 90},
					{WorkoutRecordID: r2.ID, SetNo: 1, Reps: 6, ExerciseWeight: 110},
					{WorkoutRecordID: r3.ID, SetNo: 1, Reps: 5, ExerciseWeight: 120},
				}).Error)

				return u.ID, ex1.ID
			},
			wantLen:     3,
			expectError: false,
		},
		{
			name: "【正常系】該当が無ければ空配列を返すこと",
			prepare: func(db *gorm.DB) (uint, uint) {
				u := models.User{Email: "emptysets@example.com"}
				ex := models.Exercise{Name: "ベンチ"}
				require.NoError(t, db.Create(&u).Error)
				require.NoError(t, db.Create(&ex).Error)
				return u.ID, ex.ID
			},
			wantLen:     0,
			expectError: false,
		},
		{
			name: "【異常系】workout_sets をドロップした場合はエラーを返すこと",
			prepare: func(db *gorm.DB) (uint, uint) {
				u := models.User{Email: "errsets@example.com"}
				ex := models.Exercise{Name: "プルアップ"}
				require.NoError(t, db.Create(&u).Error)
				require.NoError(t, db.Create(&ex).Error)
				require.NoError(t, db.Migrator().DropTable(&models.WorkoutSet{}))
				return u.ID, ex.ID
			},
			wantLen:     0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db := newWorkoutTestDB(t)
			repo := NewWorkoutRepository(db)

			userID, ex1ID := tt.prepare(db)

			if tt.userID == 0 {
				tt.userID = userID
			}
			if tt.exerciseID == 0 {
				tt.exerciseID = ex1ID
			}

			got, err := repo.FindSetsByUserAndExercise(tt.userID, tt.exerciseID)
			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Len(t, got, tt.wantLen)

			if tt.wantLen == 3 {
				require.Equal(t, 1, got[0].SetNo)
				require.Equal(t, 2, got[1].SetNo)
				require.Equal(t, 1, got[2].SetNo)
				require.True(t, got[0].TrainedOn.Before(got[2].TrainedOn) || got[0].TrainedOn.Equal(got[2].TrainedOn))
			}
		})
	}
}
