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

func newTimelineTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(&models.User{}, &models.Exercise{}, &models.WorkoutRecord{}))

	return db
}

func TestTimelineRepository_FindPublicRecords(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		prepare     func(db *gorm.DB)
		wantLen     int
		expectError bool
	}{
		{
			name: "【正常系】公開フラグが true の記録のみ取得できること",
			prepare: func(db *gorm.DB) {
				// ユーザー
				user := models.User{Email: "user@example.com"}
				require.NoError(t, db.Create(&user).Error)

				// 種目
				ex := models.Exercise{Name: "ベンチプレス"}
				require.NoError(t, db.Create(&ex).Error)

				// レコード（公開と非公開混在）
				records := []models.WorkoutRecord{
					{
						UserID:     user.ID,
						ExerciseID: ex.ID,
						BodyWeight: 70.5,
						TrainedOn:  now.Add(-2 * time.Hour),
						IsPublic:   true,
						Comment:    "公開1件目",
					},
					{
						UserID:     user.ID,
						ExerciseID: ex.ID,
						BodyWeight: 71.0,
						TrainedOn:  now.Add(-1 * time.Hour),
						IsPublic:   false,
						Comment:    "非公開",
					},
					{
						UserID:     user.ID,
						ExerciseID: ex.ID,
						BodyWeight: 69.8,
						TrainedOn:  now,
						IsPublic:   true,
						Comment:    "公開2件目",
					},
				}
				require.NoError(t, db.Create(&records).Error)
			},
			wantLen:     2,
			expectError: false,
		},
		{
			name: "【正常系】公開レコードが存在しない場合は空スライスを返すこと",
			prepare: func(db *gorm.DB) {
				user := models.User{Email: "user2@example.com"}
				require.NoError(t, db.Create(&user).Error)

				ex := models.Exercise{Name: "スクワット"}
				require.NoError(t, db.Create(&ex).Error)

				records := []models.WorkoutRecord{
					{
						UserID:     user.ID,
						ExerciseID: ex.ID,
						BodyWeight: 60.0,
						TrainedOn:  now,
						IsPublic:   false,
						Comment:    "非公開のみ",
					},
				}
				require.NoError(t, db.Create(&records).Error)
			},
			wantLen:     0,
			expectError: false,
		},
		{
			name: "【異常系】テーブルが存在しない場合はエラーを返すこと",
			prepare: func(db *gorm.DB) {
				require.NoError(t, db.Migrator().DropTable(&models.WorkoutRecord{}))
			},
			wantLen:     0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db := newTimelineTestDB(t)
			tt.prepare(db)

			repo := NewTimelineRepository(db)
			rows, err := repo.FindPublicRecords()

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, rows, tt.wantLen)

			if tt.wantLen == 0 {
				return
			}

			for _, r := range rows {
				require.NotZero(t, r.RecordID)
				require.NotZero(t, r.UserID)
				require.NotEmpty(t, r.UserEmail)
				require.NotEmpty(t, r.ExerciseName)
				require.InDelta(t, 0.0, r.BodyWeight, 100.0)

				require.False(t, r.TrainedOn.IsZero())
			}

			if tt.wantLen >= 2 {
				first := rows[0]
				second := rows[1]

				require.True(t,
					first.TrainedOn.After(second.TrainedOn) ||
						first.TrainedOn.Equal(second.TrainedOn),
					"rows should be ordered by trained_on DESC",
				)
			}
		})
	}
}
