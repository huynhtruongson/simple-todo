// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	pgconn "github.com/jackc/pgx/v5/pgconn"

	pgx "github.com/jackc/pgx/v5"
)

// QueryExecer is an autogenerated mock type for the QueryExecer type
type QueryExecer struct {
	mock.Mock
}

type QueryExecer_Expecter struct {
	mock *mock.Mock
}

func (_m *QueryExecer) EXPECT() *QueryExecer_Expecter {
	return &QueryExecer_Expecter{mock: &_m.Mock}
}

// Exec provides a mock function with given fields: ctx, sql, args
func (_m *QueryExecer) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, sql)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 pgconn.CommandTag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) (pgconn.CommandTag, error)); ok {
		return rf(ctx, sql, args...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) pgconn.CommandTag); ok {
		r0 = rf(ctx, sql, args...)
	} else {
		r0 = ret.Get(0).(pgconn.CommandTag)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(ctx, sql, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryExecer_Exec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exec'
type QueryExecer_Exec_Call struct {
	*mock.Call
}

// Exec is a helper method to define mock.On call
//   - ctx context.Context
//   - sql string
//   - args ...interface{}
func (_e *QueryExecer_Expecter) Exec(ctx interface{}, sql interface{}, args ...interface{}) *QueryExecer_Exec_Call {
	return &QueryExecer_Exec_Call{Call: _e.mock.On("Exec",
		append([]interface{}{ctx, sql}, args...)...)}
}

func (_c *QueryExecer_Exec_Call) Run(run func(ctx context.Context, sql string, args ...interface{})) *QueryExecer_Exec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *QueryExecer_Exec_Call) Return(_a0 pgconn.CommandTag, _a1 error) *QueryExecer_Exec_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *QueryExecer_Exec_Call) RunAndReturn(run func(context.Context, string, ...interface{}) (pgconn.CommandTag, error)) *QueryExecer_Exec_Call {
	_c.Call.Return(run)
	return _c
}

// Query provides a mock function with given fields: ctx, query, args
func (_m *QueryExecer) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 pgx.Rows
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) (pgx.Rows, error)); ok {
		return rf(ctx, query, args...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) pgx.Rows); ok {
		r0 = rf(ctx, query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pgx.Rows)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(ctx, query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryExecer_Query_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Query'
type QueryExecer_Query_Call struct {
	*mock.Call
}

// Query is a helper method to define mock.On call
//   - ctx context.Context
//   - query string
//   - args ...interface{}
func (_e *QueryExecer_Expecter) Query(ctx interface{}, query interface{}, args ...interface{}) *QueryExecer_Query_Call {
	return &QueryExecer_Query_Call{Call: _e.mock.On("Query",
		append([]interface{}{ctx, query}, args...)...)}
}

func (_c *QueryExecer_Query_Call) Run(run func(ctx context.Context, query string, args ...interface{})) *QueryExecer_Query_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *QueryExecer_Query_Call) Return(_a0 pgx.Rows, _a1 error) *QueryExecer_Query_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *QueryExecer_Query_Call) RunAndReturn(run func(context.Context, string, ...interface{}) (pgx.Rows, error)) *QueryExecer_Query_Call {
	_c.Call.Return(run)
	return _c
}

// QueryRow provides a mock function with given fields: ctx, query, args
func (_m *QueryExecer) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	var _ca []interface{}
	_ca = append(_ca, ctx, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 pgx.Row
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) pgx.Row); ok {
		r0 = rf(ctx, query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pgx.Row)
		}
	}

	return r0
}

// QueryExecer_QueryRow_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryRow'
type QueryExecer_QueryRow_Call struct {
	*mock.Call
}

// QueryRow is a helper method to define mock.On call
//   - ctx context.Context
//   - query string
//   - args ...interface{}
func (_e *QueryExecer_Expecter) QueryRow(ctx interface{}, query interface{}, args ...interface{}) *QueryExecer_QueryRow_Call {
	return &QueryExecer_QueryRow_Call{Call: _e.mock.On("QueryRow",
		append([]interface{}{ctx, query}, args...)...)}
}

func (_c *QueryExecer_QueryRow_Call) Run(run func(ctx context.Context, query string, args ...interface{})) *QueryExecer_QueryRow_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *QueryExecer_QueryRow_Call) Return(_a0 pgx.Row) *QueryExecer_QueryRow_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *QueryExecer_QueryRow_Call) RunAndReturn(run func(context.Context, string, ...interface{}) pgx.Row) *QueryExecer_QueryRow_Call {
	_c.Call.Return(run)
	return _c
}

// NewQueryExecer creates a new instance of QueryExecer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewQueryExecer(t interface {
	mock.TestingT
	Cleanup(func())
}) *QueryExecer {
	mock := &QueryExecer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
