package user_mock

import (
	"context"

	"github.com/sondev/todo-list/lib"
	user_entity "github.com/sondev/todo-list/services/user/entity"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (u *MockUserRepo) GetUsersByUsername(ctx context.Context, db lib.QueryExecer, username string) ([]user_entity.User, error) {
	args := u.Called(ctx, db, username)
	return args.Get(0).([]user_entity.User), args.Error(1)
}
func (u *MockUserRepo) CreateUser(ctx context.Context, db lib.QueryExecer, user user_entity.User) (int, error) {
	args := u.Called(ctx, db, user)
	return args.Int(0), args.Error(1)
}
func (u *MockUserRepo) GetUsersByUserIds(ctx context.Context, db lib.QueryExecer, ids []int) ([]user_entity.User, error) {
	args := u.Called(ctx, db, ids)
	return args.Get(0).([]user_entity.User), args.Error(1)
}
