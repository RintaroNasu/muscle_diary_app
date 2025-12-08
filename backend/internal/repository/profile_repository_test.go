package repository

import (
	"fmt"
	"testing"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/RintaroNasu/muscle_diary_app/utils"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newProfileTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}))
	return db
}

func TestProfileRepository_GetProfile(t *testing.T) {
	tests := []struct {
		name           string
		prepare        func(db *gorm.DB)
		userID         uint
		wantEmail      string
		expectError    bool
		expectNotFound bool
	}{
		{
			name: "【正常系】存在するユーザーを取得できること",
			prepare: func(db *gorm.DB) {
				user := models.User{
					Email:      "test@example.com",
					Height:     utils.Ptr(170.5),
					GoalWeight: utils.Ptr(60.0),
				}
				require.NoError(t, db.Create(&user).Error)
			},
			userID:         1,
			wantEmail:      "test@example.com",
			expectError:    false,
			expectNotFound: false,
		},
		{
			name: "【異常系】存在しないユーザーIDを指定した場合、ErrNotFoundを返すこと",
			prepare: func(db *gorm.DB) {
				// ユーザーは作成しない（=存在しないIDを指定）
			},
			userID:         99,
			wantEmail:      "",
			expectError:    true,
			expectNotFound: true,
		},
		{
			name: "【異常系】テーブルが存在しない場合、ErrNotFound以外のエラーを返すこと",
			prepare: func(db *gorm.DB) {
				require.NoError(t, db.Migrator().DropTable(&models.User{}))
			},
			userID:         1,
			wantEmail:      "",
			expectError:    true,
			expectNotFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newProfileTestDB(t)
			tt.prepare(db)

			repo := NewProfileRepository(db)
			got, err := repo.GetProfile(tt.userID)

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, got)

				if tt.expectNotFound {
					require.ErrorIs(t, err, ErrNotFound)
				} else {
					require.NotErrorIs(t, err, ErrNotFound)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, tt.wantEmail, got.Email)
		})
	}
}

func TestProfileRepository_UpdateProfile(t *testing.T) {
	tests := []struct {
		name           string
		prepare        func(db *gorm.DB)
		userID         uint
		newHeight      float64
		newGoal        float64
		expectError    bool
		expectNotFound bool
	}{
		{
			name: "【正常系】身長と目標体重を更新できること",
			prepare: func(db *gorm.DB) {
				user := models.User{
					Email:      "test@example.com",
					Height:     utils.Ptr(160.0),
					GoalWeight: utils.Ptr(55.0),
				}
				require.NoError(t, db.Create(&user).Error)
			},
			userID:         1,
			newHeight:      170.0,
			newGoal:        60.0,
			expectError:    false,
			expectNotFound: false,
		},
		{
			name: "【異常系】存在しないユーザーIDを更新した場合、ErrNotFoundを返すこと",
			prepare: func(db *gorm.DB) {
				user := models.User{
					Email: "dummy@example.com",
				}
				require.NoError(t, db.Create(&user).Error)
			},
			userID:         99,
			newHeight:      180.0,
			newGoal:        70.0,
			expectError:    true,
			expectNotFound: true,
		},
		{
			name: "【異常系】テーブルが存在しない場合、エラーを返すこと",
			prepare: func(db *gorm.DB) {
				require.NoError(t, db.Migrator().DropTable(&models.User{}))
			},
			userID:         1,
			newHeight:      175.0,
			newGoal:        65.0,
			expectError:    true,
			expectNotFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newProfileTestDB(t)
			tt.prepare(db)

			repo := NewProfileRepository(db)
			err := repo.UpdateProfile(tt.userID, &tt.newHeight, &tt.newGoal)

			if tt.expectError {
				require.Error(t, err)

				if tt.expectNotFound {
					require.ErrorIs(t, err, ErrNotFound)
				}
				return
			}

			require.NoError(t, err)

			var user models.User
			require.NoError(t, db.First(&user, tt.userID).Error)
			require.Equal(t, tt.newHeight, *user.Height)
			require.Equal(t, tt.newGoal, *user.GoalWeight)
		})
	}
}
