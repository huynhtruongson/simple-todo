package user_service

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

type MockServiceProp struct {
	DB       *mock_db.DB
	TX       *mock_db.Tx
	UserRepo *mock_repo.UserRepo
}

func NewMockUserService(t *testing.T) (*UserService, *MockServiceProp) {
	userRepo := mock_repo.NewUserRepo(t)
	db := mock_db.NewDB(t)
	tx := mock_db.NewTx(t)
	return &UserService{
			DB:       db,
			UserRepo: userRepo,
		}, &MockServiceProp{
			DB:       db,
			TX:       tx,
			UserRepo: userRepo,
		}
}
func TestCreateUserBiz_CreateUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name         string
		user         user_entity.User
		mock         func(prop *MockServiceProp)
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
			mock: func(prop *MockServiceProp) {
				prop.UserRepo.EXPECT().GetUsersByUsername(ctx, prop.DB, "username").Once().Return([]user_entity.User{}, nil)
				prop.DB.EXPECT().BeginTx(ctx, mock.Anything).Once().Return(prop.TX, nil)
				prop.UserRepo.EXPECT().CreateUser(ctx, prop.TX, mock.Anything).Once().Return(1, nil)
				prop.TX.EXPECT().Commit(ctx).Once().Return(nil)
			},
			expectErr:    nil,
			expectUserID: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sv, prop := NewMockUserService(t)
			tt.mock(prop)
			userID, err := sv.CreateUser(ctx, tt.user)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
			assert.Equal(t, tt.expectUserID, userID)
		})
	}
}

func TestCreateUserBiz_ValidateUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tests := []struct {
		name      string
		user      user_entity.User
		mock      func(prop *MockServiceProp)
		expectErr *common.AppError
	}{
		{
			name: "should throw error when fullname is empty",
			user: user_entity.User{
				FullName: "",
			},
			mock:      func(prop *MockServiceProp) {},
			expectErr: common.NewInvalidRequestError(nil, user_entity.ErrorFullnameIsEmpty, "ValidateUser"),
		},
		{
			name: "should throw error when username is empty",
			user: user_entity.User{
				FullName: "fullname",
				Username: "",
			},
			mock:      func(prop *MockServiceProp) {},
			expectErr: common.NewInvalidRequestError(nil, user_entity.ErrorUsernameIsEmpty, "ValidateUser"),
		},
		{
			name: "should throw error when username is less than 6 characters",
			user: user_entity.User{
				FullName: "fullname",
				Username: "user",
			},
			mock:      func(prop *MockServiceProp) {},
			expectErr: common.NewInvalidRequestError(nil, user_entity.ErrorInvalidUsernameLength, "ValidateUser"),
		},
		{
			name: "should throw error when password is empty",
			user: user_entity.User{
				FullName: "fullname",
				Username: "username",
				Password: "",
			},
			mock:      func(prop *MockServiceProp) {},
			expectErr: common.NewInvalidRequestError(nil, user_entity.ErrorPasswordIsEmpty, "ValidateUser"),
		},
		{
			name: "should throw error when password is less than 6 characters",
			user: user_entity.User{
				FullName: "fullname",
				Username: "username",
				Password: "123",
			},
			mock:      func(prop *MockServiceProp) {},
			expectErr: common.NewInvalidRequestError(nil, user_entity.ErrorInvalidPasswordLength, "ValidateUser"),
		},
		{
			name: "should throw error when username has already existed",
			user: user_entity.User{
				FullName: "fullname",
				Username: "username",
				Password: "123123",
			},
			mock: func(prop *MockServiceProp) {
				prop.UserRepo.EXPECT().GetUsersByUsername(ctx, prop.DB, "username").Once().Return([]user_entity.User{{
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
			sv, prop := NewMockUserService(t)
			tt.mock(prop)
			err := sv.ValidateUser(ctx, tt.user)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr.Code, err.(*common.AppError).Code)
				assert.Equal(t, tt.expectErr.Message, err.(*common.AppError).Message)
			}
		})
	}
}
