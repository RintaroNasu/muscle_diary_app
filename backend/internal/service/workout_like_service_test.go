package service

import (
	"errors"
	"testing"

	"github.com/RintaroNasu/muscle_diary_app/internal/repository"
	"github.com/stretchr/testify/require"
)

type fakeWorkoutLikeRepo struct {
	isRecordPublicFunc func(recordID uint) (bool, error)
	createLikeFunc     func(userID uint, recordID uint) error
	deleteLikeFunc     func(userID uint, recordID uint) error
	isLikedByMeFunc    func(userID uint, recordID uint) (bool, error)

	createCalled int
	deleteCalled int
}

func (f *fakeWorkoutLikeRepo) CreateLike(userID uint, recordID uint) error {
	f.createCalled++
	if f.createLikeFunc == nil {
		return nil
	}
	return f.createLikeFunc(userID, recordID)
}

func (f *fakeWorkoutLikeRepo) DeleteLike(userID uint, recordID uint) error {
	f.deleteCalled++
	if f.deleteLikeFunc == nil {
		return nil
	}
	return f.deleteLikeFunc(userID, recordID)
}

func (f *fakeWorkoutLikeRepo) IsRecordPublic(recordID uint) (bool, error) {
	if f.isRecordPublicFunc == nil {
		return true, nil
	}
	return f.isRecordPublicFunc(recordID)
}

func (f *fakeWorkoutLikeRepo) IsLikedByMe(userID uint, recordID uint) (bool, error) {
	if f.isLikedByMeFunc == nil {
		return false, nil
	}
	return f.isLikedByMeFunc(userID, recordID)
}

func TestWorkoutLikeService_Like(t *testing.T) {
	tests := []struct {
		name        string
		repo        fakeWorkoutLikeRepo
		userID      uint
		recordID    uint
		wantErr     error
		errContains string
		wantCreate  int
		wantDelete  int
	}{
		{
			name:   "【正常系】公開レコードならCreateLikeが呼ばれて成功する",
			userID: 1, recordID: 10,
			repo: fakeWorkoutLikeRepo{
				isRecordPublicFunc: func(recordID uint) (bool, error) {
					require.Equal(t, uint(10), recordID)
					return true, nil
				},
				createLikeFunc: func(userID uint, recordID uint) error {
					require.Equal(t, uint(1), userID)
					require.Equal(t, uint(10), recordID)
					return nil
				},
			},
			wantCreate: 1,
		},
		{
			name:   "【正常系】非公開レコードならErrForbiddenPrivateRecordを返しCreateLikeは呼ばれない",
			userID: 1, recordID: 10,
			repo: fakeWorkoutLikeRepo{
				isRecordPublicFunc: func(recordID uint) (bool, error) {
					return false, nil
				},
			},
			wantErr:    ErrForbiddenPrivateRecord,
			wantCreate: 0,
		},
		{
			name:   "【正常系】存在しないレコードならErrRecordNotFoundに変換して返す",
			userID: 1, recordID: 999,
			repo: fakeWorkoutLikeRepo{
				isRecordPublicFunc: func(recordID uint) (bool, error) {
					return false, repository.ErrNotFound
				},
			},
			wantErr:    ErrRecordNotFound,
			wantCreate: 0,
		},
		{
			name:   "【異常系】IsRecordPublicが想定外エラーならそのまま返す",
			userID: 1, recordID: 10,
			repo: fakeWorkoutLikeRepo{
				isRecordPublicFunc: func(recordID uint) (bool, error) {
					return false, errors.New("db down")
				},
			},
			errContains: "db down",
			wantCreate:  0,
		},
		{
			name:   "【異常系】CreateLikeがエラーならそのまま返す",
			userID: 1, recordID: 10,
			repo: fakeWorkoutLikeRepo{
				isRecordPublicFunc: func(recordID uint) (bool, error) {
					return true, nil
				},
				createLikeFunc: func(userID uint, recordID uint) error {
					return errors.New("insert failed")
				},
			},
			errContains: "insert failed",
			wantCreate:  1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.repo
			svc := NewWorkoutLikeService(&repo)

			err := svc.Like(tt.userID, tt.recordID)

			switch {
			case tt.wantErr != nil:
				require.ErrorIs(t, err, tt.wantErr)
			case tt.errContains != "":
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errContains)
			default:
				require.NoError(t, err)
			}

			require.Equal(t, tt.wantCreate, repo.createCalled)
			require.Equal(t, tt.wantDelete, repo.deleteCalled)
		})
	}
}

func TestWorkoutLikeService_Unlike(t *testing.T) {
	tests := []struct {
		name        string
		repo        fakeWorkoutLikeRepo
		userID      uint
		recordID    uint
		wantErr     error
		errContains string
		wantCreate  int
		wantDelete  int
	}{
		{
			name:   "【正常系】公開レコードならDeleteLikeが呼ばれて成功する（冪等）",
			userID: 1, recordID: 10,
			repo: fakeWorkoutLikeRepo{
				isRecordPublicFunc: func(recordID uint) (bool, error) {
					return true, nil
				},
				deleteLikeFunc: func(userID uint, recordID uint) error {
					require.Equal(t, uint(1), userID)
					require.Equal(t, uint(10), recordID)
					return nil
				},
			},
			wantDelete: 1,
		},
		{
			name:   "【正常系】非公開レコードならErrForbiddenPrivateRecordを返しDeleteLikeは呼ばれない",
			userID: 1, recordID: 10,
			repo: fakeWorkoutLikeRepo{
				isRecordPublicFunc: func(recordID uint) (bool, error) {
					return false, nil
				},
			},
			wantErr:    ErrForbiddenPrivateRecord,
			wantDelete: 0,
		},
		{
			name:   "【正常系】存在しないレコードならErrRecordNotFoundに変換して返す",
			userID: 1, recordID: 999,
			repo: fakeWorkoutLikeRepo{
				isRecordPublicFunc: func(recordID uint) (bool, error) {
					return false, repository.ErrNotFound
				},
			},
			wantErr:    ErrRecordNotFound,
			wantDelete: 0,
		},
		{
			name:   "【異常系】IsRecordPublicが想定外エラーならそのまま返す",
			userID: 1, recordID: 10,
			repo: fakeWorkoutLikeRepo{
				isRecordPublicFunc: func(recordID uint) (bool, error) {
					return false, errors.New("db down")
				},
			},
			errContains: "db down",
			wantDelete:  0,
		},
		{
			name:   "【異常系】DeleteLikeがエラーならそのまま返す",
			userID: 1, recordID: 10,
			repo: fakeWorkoutLikeRepo{
				isRecordPublicFunc: func(recordID uint) (bool, error) {
					return true, nil
				},
				deleteLikeFunc: func(userID uint, recordID uint) error {
					return errors.New("delete failed")
				},
			},
			errContains: "delete failed",
			wantDelete:  1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.repo
			svc := NewWorkoutLikeService(&repo)

			err := svc.Unlike(tt.userID, tt.recordID)

			switch {
			case tt.wantErr != nil:
				require.ErrorIs(t, err, tt.wantErr)
			case tt.errContains != "":
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errContains)
			default:
				require.NoError(t, err)
			}

			require.Equal(t, tt.wantCreate, repo.createCalled)
			require.Equal(t, tt.wantDelete, repo.deleteCalled)
		})
	}
}
