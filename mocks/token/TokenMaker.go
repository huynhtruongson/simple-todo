// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"

	token "github.com/huynhtruongson/simple-todo/token"
)

// TokenMaker is an autogenerated mock type for the TokenMaker type
type TokenMaker struct {
	mock.Mock
}

type TokenMaker_Expecter struct {
	mock *mock.Mock
}

func (_m *TokenMaker) EXPECT() *TokenMaker_Expecter {
	return &TokenMaker_Expecter{mock: &_m.Mock}
}

// CreateToken provides a mock function with given fields: userID, duration, tokenType
func (_m *TokenMaker) CreateToken(userID int, duration time.Duration, tokenType token.TokenType) (string, token.TokenPayload, error) {
	ret := _m.Called(userID, duration, tokenType)

	var r0 string
	var r1 token.TokenPayload
	var r2 error
	if rf, ok := ret.Get(0).(func(int, time.Duration, token.TokenType) (string, token.TokenPayload, error)); ok {
		return rf(userID, duration, tokenType)
	}
	if rf, ok := ret.Get(0).(func(int, time.Duration, token.TokenType) string); ok {
		r0 = rf(userID, duration, tokenType)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(int, time.Duration, token.TokenType) token.TokenPayload); ok {
		r1 = rf(userID, duration, tokenType)
	} else {
		r1 = ret.Get(1).(token.TokenPayload)
	}

	if rf, ok := ret.Get(2).(func(int, time.Duration, token.TokenType) error); ok {
		r2 = rf(userID, duration, tokenType)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// TokenMaker_CreateToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateToken'
type TokenMaker_CreateToken_Call struct {
	*mock.Call
}

// CreateToken is a helper method to define mock.On call
//   - userID int
//   - duration time.Duration
//   - tokenType token.TokenType
func (_e *TokenMaker_Expecter) CreateToken(userID interface{}, duration interface{}, tokenType interface{}) *TokenMaker_CreateToken_Call {
	return &TokenMaker_CreateToken_Call{Call: _e.mock.On("CreateToken", userID, duration, tokenType)}
}

func (_c *TokenMaker_CreateToken_Call) Run(run func(userID int, duration time.Duration, tokenType token.TokenType)) *TokenMaker_CreateToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int), args[1].(time.Duration), args[2].(token.TokenType))
	})
	return _c
}

func (_c *TokenMaker_CreateToken_Call) Return(_a0 string, _a1 token.TokenPayload, _a2 error) *TokenMaker_CreateToken_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *TokenMaker_CreateToken_Call) RunAndReturn(run func(int, time.Duration, token.TokenType) (string, token.TokenPayload, error)) *TokenMaker_CreateToken_Call {
	_c.Call.Return(run)
	return _c
}

// VerifyToken provides a mock function with given fields: _a0
func (_m *TokenMaker) VerifyToken(_a0 string) (token.TokenPayload, error) {
	ret := _m.Called(_a0)

	var r0 token.TokenPayload
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (token.TokenPayload, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) token.TokenPayload); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(token.TokenPayload)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TokenMaker_VerifyToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'VerifyToken'
type TokenMaker_VerifyToken_Call struct {
	*mock.Call
}

// VerifyToken is a helper method to define mock.On call
//   - _a0 string
func (_e *TokenMaker_Expecter) VerifyToken(_a0 interface{}) *TokenMaker_VerifyToken_Call {
	return &TokenMaker_VerifyToken_Call{Call: _e.mock.On("VerifyToken", _a0)}
}

func (_c *TokenMaker_VerifyToken_Call) Run(run func(_a0 string)) *TokenMaker_VerifyToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *TokenMaker_VerifyToken_Call) Return(_a0 token.TokenPayload, _a1 error) *TokenMaker_VerifyToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TokenMaker_VerifyToken_Call) RunAndReturn(run func(string) (token.TokenPayload, error)) *TokenMaker_VerifyToken_Call {
	_c.Call.Return(run)
	return _c
}

// NewTokenMaker creates a new instance of TokenMaker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTokenMaker(t interface {
	mock.TestingT
	Cleanup(func())
}) *TokenMaker {
	mock := &TokenMaker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
