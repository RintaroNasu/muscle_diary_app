package repository

import (
	"testing"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}))

	return db
}

func TestUserRepository_CreateAndFindByEmail(t *testing.T) {
	db := newTestDB(t)
	repo := NewUserRepository(db)

	t.Run("【正常系】ユーザーを作成できること", func(t *testing.T) {
		u := &models.User{Email: "test@test.com", Password: "asdfasdf"}
		err := repo.Create(u)
		require.NoError(t, err)
		require.NotZero(t, u.ID)

		got, err := repo.FindByEmail("test@test.com")
		require.NoError(t, err)
		require.Equal(t, u.ID, got.ID)
		require.Equal(t, "test@test.com", got.Email)
	})

	t.Run("【正常系】存在しないメールアドレスを検索するとErrNotFoundを返すこと", func(t *testing.T) {
		got, err := repo.FindByEmail("none@test.com")
		require.ErrorIs(t, err, ErrNotFound)
		require.Nil(t, got)
	})

	t.Run("【異常系】同じメールアドレスを2回登録するとErrUniqueViolationを返すこと", func(t *testing.T) {
		u1 := &models.User{Email: "dup@test.com", Password: "asdfasdf"}
		require.NoError(t, repo.Create(u1))

		u2 := &models.User{Email: "dup@test.com", Password: "asdfasdf"}
		err := repo.Create(u2)
		require.ErrorIs(t, err, ErrUniqueViolation)
	})
}
