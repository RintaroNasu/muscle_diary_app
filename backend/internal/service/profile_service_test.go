package service

import (
	"errors"
	"testing"

	"github.com/RintaroNasu/muscle_diary_app/internal/models"
	"github.com/RintaroNasu/muscle_diary_app/utils"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type fakeProfileRepo struct {
	getFunc    func(userID uint) (*models.User, error)
	updateFunc func(userID uint, height *float64, goalWeight *float64) error
}

func (f *fakeProfileRepo) GetProfile(userID uint) (*models.User, error) {
	return f.getFunc(userID)
}

func (f *fakeProfileRepo) UpdateProfile(userID uint, height *float64, goalWeight *float64) error {
	return f.updateFunc(userID, height, goalWeight)
}

func TestProfileService_GetProfile(t *testing.T) {
	tests := []struct {
		name        string
		mockRepo    fakeProfileRepo
		userID      uint
		wantUser    *models.User
		wantErr     error
		errContains string
	}{
		{
			name: "【正常系】ユーザー情報を取得できること",
			mockRepo: fakeProfileRepo{
				getFunc: func(userID uint) (*models.User, error) {
					return &models.User{
						Email:      "test@example.com",
						Height:     utils.Ptr(170.0),
						GoalWeight: utils.Ptr(60.0),
					}, nil
				},
			},
			userID:   1,
			wantUser: &models.User{Email: "test@example.com"},
		},
		{
			name: "【異常系】存在しないユーザーの場合は ErrUserNotFound を返すこと",
			mockRepo: fakeProfileRepo{
				getFunc: func(userID uint) (*models.User, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			userID:  99,
			wantErr: ErrUserNotFound,
		},
		{
			name: "【異常系】DB障害などでエラーが発生した場合はそのまま返すこと",
			mockRepo: fakeProfileRepo{
				getFunc: func(userID uint) (*models.User, error) {
					return nil, errors.New("db down")
				},
			},
			userID:      1,
			errContains: "db down",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProfileService(&tt.mockRepo)
			got, err := svc.GetProfile(tt.userID)

			switch {
			case tt.wantErr != nil:
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, got)
			case tt.errContains != "":
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errContains)
			default:
				require.NoError(t, err)
				require.NotNil(t, got)
				require.Equal(t, tt.wantUser.Email, got.Email)
			}
		})
	}
}

func TestProfileService_UpdateProfile(t *testing.T) {
	tests := []struct {
		name        string
		mockRepo    fakeProfileRepo
		userID      uint
		height      *float64
		goalWeight  *float64
		wantUser    *models.User
		wantErr     error
		errContains string
	}{
		{
			name: "【正常系】身長・目標体重を更新して取得できること",
			mockRepo: fakeProfileRepo{
				updateFunc: func(userID uint, height *float64, goalWeight *float64) error {
					return nil
				},
				getFunc: func(userID uint) (*models.User, error) {
					return &models.User{
						Email:      "updated@example.com",
						Height:     utils.Ptr(175.0),
						GoalWeight: utils.Ptr(65.0),
					}, nil
				},
			},
			userID:     1,
			height:     utils.Ptr(175.0),
			goalWeight: utils.Ptr(65.0),
			wantUser:   &models.User{Email: "updated@example.com"},
		},
		{
			name: "【異常系】UpdateProfileでエラーが発生した場合はそのまま返すこと",
			mockRepo: fakeProfileRepo{
				updateFunc: func(userID uint, height *float64, goalWeight *float64) error {
					return errors.New("update failed")
				},
			},
			userID:      1,
			height:      utils.Ptr(180.0),
			goalWeight:  utils.Ptr(70.0),
			errContains: "update failed",
		},
		{
			name: "【異常系】更新後のGetProfileでエラーが発生した場合はそのまま返すこと",
			mockRepo: fakeProfileRepo{
				updateFunc: func(userID uint, height *float64, goalWeight *float64) error {
					return nil
				},
				getFunc: func(userID uint) (*models.User, error) {
					return nil, errors.New("select failed")
				},
			},
			userID:      1,
			height:      utils.Ptr(180.0),
			goalWeight:  utils.Ptr(70.0),
			errContains: "select failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProfileService(&tt.mockRepo)
			got, err := svc.UpdateProfile(tt.userID, tt.height, tt.goalWeight)

			switch {
			case tt.wantErr != nil:
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, got)
			case tt.errContains != "":
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errContains)
			default:
				require.NoError(t, err)
				require.NotNil(t, got)
				require.Equal(t, tt.wantUser.Email, got.Email)
				require.Equal(t, *tt.height, *got.Height)
				require.Equal(t, *tt.goalWeight, *got.GoalWeight)
			}
		})
	}
}
