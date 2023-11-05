package user_biz

import (
	"context"
	"testing"

	"github.com/huynhtruongson/simple-todo/common"
	mock_db "github.com/huynhtruongson/simple-todo/mocks/lib"
	mock_repo "github.com/huynhtruongson/simple-todo/mocks/user"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUserBiz_CreateUser(t *testing.T) {
	ctx := context.Background()
	userRepo := mock_repo.NewUserRepo(t)
	db := mock_db.NewDB(t)
	tx := mock_db.NewTx(t)

	tests := []struct {
		name         string
		user         user_entity.User
		mock         func()
		expectErr    error
		expectUserID int
	}{
		{
			name: "should return userID when create user successfully",
			user: user_entity.User{
				FullName: "fullname",
				Username: "username",
				Password: "123123",
			},
			mock: func() {
				// userRepo.On("GetUsersByUsername", ctx, db, "username").Once().Return([]user_entity.User{}, nil)
				userRepo.EXPECT().GetUsersByUsername(ctx, db, "username").Once().Return([]user_entity.User{}, nil)
				db.EXPECT().BeginTx(ctx, mock.Anything).Once().Return(tx, nil)
				userRepo.EXPECT().CreateUser(ctx, tx, user_entity.User{
					FullName: "fullname",
					Username: "username",
					Password: "123123",
				}).Once().Return(1, nil)
				tx.EXPECT().Commit(ctx).Once().Return(nil)
			},
			expectErr:    nil,
			expectUserID: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			biz := NewCreateUserBiz(db, userRepo)
			tt.mock()
			userID, err := biz.CreateUser(ctx, tt.user)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
			assert.Equal(t, tt.expectUserID, userID)
		})
	}
}

func TestCreateUserBiz_ValidateUser(t *testing.T) {
	ctx := context.Background()
	userRepo := mock_repo.NewUserRepo(t)
	db := mock_db.NewDB(t)
	tests := []struct {
		name      string
		user      user_entity.User
		mock      func()
		expectErr *common.AppError
	}{
		{
			name: "should throw error when fullname is empty",
			user: user_entity.User{
				FullName: "",
			},
			mock:      func() {},
			expectErr: common.NewInvalidRequestError(nil, user_entity.ErrorFullnameIsEmpty, "ValidateUser"),
		},
		{
			name: "should throw error when username is empty",
			user: user_entity.User{
				FullName: "fullname",
				Username: "",
			},
			mock:      func() {},
			expectErr: common.NewInvalidRequestError(nil, user_entity.ErrorUsernameIsEmpty, "ValidateUser"),
		},
		{
			name: "should throw error when username is less than 6 characters",
			user: user_entity.User{
				FullName: "fullname",
				Username: "user",
			},
			mock:      func() {},
			expectErr: common.NewInvalidRequestError(nil, user_entity.ErrorInvalidUsernameLength, "ValidateUser"),
		},
		{
			name: "should throw error when password is empty",
			user: user_entity.User{
				FullName: "fullname",
				Username: "username",
				Password: "",
			},
			mock:      func() {},
			expectErr: common.NewInvalidRequestError(nil, user_entity.ErrorPasswordIsEmpty, "ValidateUser"),
		},
		{
			name: "should throw error when password is less than 6 characters",
			user: user_entity.User{
				FullName: "fullname",
				Username: "username",
				Password: "123",
			},
			mock:      func() {},
			expectErr: common.NewInvalidRequestError(nil, user_entity.ErrorInvalidPasswordLength, "ValidateUser"),
		},
		{
			name: "should throw error when username has already existed",
			user: user_entity.User{
				FullName: "fullname",
				Username: "username",
				Password: "123123",
			},
			mock: func() {
				userRepo.EXPECT().GetUsersByUsername(ctx, db, "username").Once().Return([]user_entity.User{{
					FullName: "fullname",
					Username: "username",
					Password: "123123",
				}}, nil)
			},
			expectErr: common.NewInvalidRequestError(nil, user_entity.ErrorUserNameAlreadyExist, "ValidateUser"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			biz := CreateUserBiz{
				DB:       db,
				UserRepo: userRepo,
			}
			tt.mock()
			err := biz.ValidateUser(ctx, tt.user)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr.Code, err.(*common.AppError).Code)
				assert.Equal(t, tt.expectErr.Message, err.(*common.AppError).Message)
			}
		})
	}
}
