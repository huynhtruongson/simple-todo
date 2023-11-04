package user_biz

import (
	"context"
	"testing"

	"github.com/sondev/todo-list/common"
	mock_db "github.com/sondev/todo-list/mock"
	user_entity "github.com/sondev/todo-list/services/user/entity"
	user_mock "github.com/sondev/todo-list/services/user/mock/repo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUserBiz_CreateUser(t *testing.T) {
	ctx := context.Background()
	userRepo := &user_mock.MockUserRepo{}
	db := &mock_db.MockDB{}
	tx := &mock_db.MockTx{}
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
				userRepo.On("GetUsersByUsername", ctx, db, "username").Once().Return([]user_entity.User{}, nil)
				db.On("BeginTx", ctx, mock.Anything).Once().Return(tx, nil)
				userRepo.On("CreateUser", ctx, tx, user_entity.User{
					FullName: "fullname",
					Username: "username",
					Password: "123123",
				}).Once().Return(1, nil)
				tx.On("Commit", ctx).Once().Return(nil)
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
	userRepo := &user_mock.MockUserRepo{}
	db := &mock_db.MockDB{}
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
				userRepo.On("GetUsersByUsername", ctx, db, "username").Once().Return([]user_entity.User{{
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
