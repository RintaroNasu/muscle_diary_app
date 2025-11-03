package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newExerciseTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.Exercise{}))
	return db
}

func TestExerciseRepository_List(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		prepare     func(db *gorm.DB)
		wantNames   []string
		wantLen     int
		expectError bool
	}{
		{
			name: "【正常系】レコードが存在する場合、id ASC で (id, name) を取得できること",
			prepare: func(db *gorm.DB) {
				seeds := []models.Exercise{
					{Name: "Bench Press"},
					{Name: "Squat"},
					{Name: "Deadlift"},
				}
				require.NoError(t, db.Create(&seeds).Error)
			},
			wantNames:   []string{"Bench Press", "Squat", "Deadlift"},
			wantLen:     3,
			expectError: false,
		},
		{
			name:        "【正常系】レコードが0件の場合、空スライスを返すこと（エラーなし）",
			prepare:     func(db *gorm.DB) {},
			wantNames:   []string{},
			wantLen:     0,
			expectError: false,
		},
		{
			name: "【異常系】テーブルが存在しない場合はエラーを返すこと",
			prepare: func(db *gorm.DB) {
				require.NoError(t, db.Migrator().DropTable(&models.Exercise{}))
			},
			wantNames:   nil,
			wantLen:     0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newExerciseTestDB(t)
			tt.prepare(db)

			repo := NewExerciseRepository(db)
			got, err := repo.List(ctx)

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Len(t, got, tt.wantLen)

			if tt.wantLen > 0 {
				names := make([]string, len(got))
				for i, e := range got {
					names[i] = e.Name
				}
				require.ElementsMatch(t, tt.wantNames, names)
			}
		})
	}
}
