// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"

	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"
	mock "github.com/stretchr/testify/mock"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

type UserService_Expecter struct {
	mock *mock.Mock
}

func (_m *UserService) EXPECT() *UserService_Expecter {
	return &UserService_Expecter{mock: &_m.Mock}
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *UserService) CreateUser(ctx context.Context, user user_entity.User) (int, error) {
	ret := _m.Called(ctx, user)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, user_entity.User) (int, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, user_entity.User) int); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, user_entity.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserService_CreateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUser'
type UserService_CreateUser_Call struct {
	*mock.Call
}

// CreateUser is a helper method to define mock.On call
//   - ctx context.Context
//   - user user_entity.User
func (_e *UserService_Expecter) CreateUser(ctx interface{}, user interface{}) *UserService_CreateUser_Call {
	return &UserService_CreateUser_Call{Call: _e.mock.On("CreateUser", ctx, user)}
}

func (_c *UserService_CreateUser_Call) Run(run func(ctx context.Context, user user_entity.User)) *UserService_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(user_entity.User))
	})
	return _c
}

func (_c *UserService_CreateUser_Call) Return(_a0 int, _a1 error) *UserService_CreateUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserService_CreateUser_Call) RunAndReturn(run func(context.Context, user_entity.User) (int, error)) *UserService_CreateUser_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserService creates a new instance of UserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserService(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserService {
	mock := &UserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
