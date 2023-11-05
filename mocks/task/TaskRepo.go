// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"

	lib "github.com/huynhtruongson/simple-todo/lib"
	mock "github.com/stretchr/testify/mock"

	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"
)

// TaskRepo is an autogenerated mock type for the TaskRepo type
type TaskRepo struct {
	mock.Mock
}

type TaskRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *TaskRepo) EXPECT() *TaskRepo_Expecter {
	return &TaskRepo_Expecter{mock: &_m.Mock}
}

// CountTask provides a mock function with given fields: ctx, db
func (_m *TaskRepo) CountTask(ctx context.Context, db lib.QueryExecer) (int, error) {
	ret := _m.Called(ctx, db)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, lib.QueryExecer) (int, error)); ok {
		return rf(ctx, db)
	}
	if rf, ok := ret.Get(0).(func(context.Context, lib.QueryExecer) int); ok {
		r0 = rf(ctx, db)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, lib.QueryExecer) error); ok {
		r1 = rf(ctx, db)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TaskRepo_CountTask_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CountTask'
type TaskRepo_CountTask_Call struct {
	*mock.Call
}

// CountTask is a helper method to define mock.On call
//   - ctx context.Context
//   - db lib.QueryExecer
func (_e *TaskRepo_Expecter) CountTask(ctx interface{}, db interface{}) *TaskRepo_CountTask_Call {
	return &TaskRepo_CountTask_Call{Call: _e.mock.On("CountTask", ctx, db)}
}

func (_c *TaskRepo_CountTask_Call) Run(run func(ctx context.Context, db lib.QueryExecer)) *TaskRepo_CountTask_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(lib.QueryExecer))
	})
	return _c
}

func (_c *TaskRepo_CountTask_Call) Return(_a0 int, _a1 error) *TaskRepo_CountTask_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TaskRepo_CountTask_Call) RunAndReturn(run func(context.Context, lib.QueryExecer) (int, error)) *TaskRepo_CountTask_Call {
	_c.Call.Return(run)
	return _c
}

// CreateTask provides a mock function with given fields: ctx, db, task
func (_m *TaskRepo) CreateTask(ctx context.Context, db lib.QueryExecer, task task_entity.Task) (int, error) {
	ret := _m.Called(ctx, db, task)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, lib.QueryExecer, task_entity.Task) (int, error)); ok {
		return rf(ctx, db, task)
	}
	if rf, ok := ret.Get(0).(func(context.Context, lib.QueryExecer, task_entity.Task) int); ok {
		r0 = rf(ctx, db, task)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, lib.QueryExecer, task_entity.Task) error); ok {
		r1 = rf(ctx, db, task)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TaskRepo_CreateTask_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateTask'
type TaskRepo_CreateTask_Call struct {
	*mock.Call
}

// CreateTask is a helper method to define mock.On call
//   - ctx context.Context
//   - db lib.QueryExecer
//   - task task_entity.Task
func (_e *TaskRepo_Expecter) CreateTask(ctx interface{}, db interface{}, task interface{}) *TaskRepo_CreateTask_Call {
	return &TaskRepo_CreateTask_Call{Call: _e.mock.On("CreateTask", ctx, db, task)}
}

func (_c *TaskRepo_CreateTask_Call) Run(run func(ctx context.Context, db lib.QueryExecer, task task_entity.Task)) *TaskRepo_CreateTask_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(lib.QueryExecer), args[2].(task_entity.Task))
	})
	return _c
}

func (_c *TaskRepo_CreateTask_Call) Return(_a0 int, _a1 error) *TaskRepo_CreateTask_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TaskRepo_CreateTask_Call) RunAndReturn(run func(context.Context, lib.QueryExecer, task_entity.Task) (int, error)) *TaskRepo_CreateTask_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteTask provides a mock function with given fields: ctx, db, id
func (_m *TaskRepo) DeleteTask(ctx context.Context, db lib.QueryExecer, id int) error {
	ret := _m.Called(ctx, db, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, lib.QueryExecer, int) error); ok {
		r0 = rf(ctx, db, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TaskRepo_DeleteTask_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteTask'
type TaskRepo_DeleteTask_Call struct {
	*mock.Call
}

// DeleteTask is a helper method to define mock.On call
//   - ctx context.Context
//   - db lib.QueryExecer
//   - id int
func (_e *TaskRepo_Expecter) DeleteTask(ctx interface{}, db interface{}, id interface{}) *TaskRepo_DeleteTask_Call {
	return &TaskRepo_DeleteTask_Call{Call: _e.mock.On("DeleteTask", ctx, db, id)}
}

func (_c *TaskRepo_DeleteTask_Call) Run(run func(ctx context.Context, db lib.QueryExecer, id int)) *TaskRepo_DeleteTask_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(lib.QueryExecer), args[2].(int))
	})
	return _c
}

func (_c *TaskRepo_DeleteTask_Call) Return(_a0 error) *TaskRepo_DeleteTask_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TaskRepo_DeleteTask_Call) RunAndReturn(run func(context.Context, lib.QueryExecer, int) error) *TaskRepo_DeleteTask_Call {
	_c.Call.Return(run)
	return _c
}

// GetTasksByIds provides a mock function with given fields: ctx, db, ids
func (_m *TaskRepo) GetTasksByIds(ctx context.Context, db lib.QueryExecer, ids []int) ([]task_entity.Task, error) {
	ret := _m.Called(ctx, db, ids)

	var r0 []task_entity.Task
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, lib.QueryExecer, []int) ([]task_entity.Task, error)); ok {
		return rf(ctx, db, ids)
	}
	if rf, ok := ret.Get(0).(func(context.Context, lib.QueryExecer, []int) []task_entity.Task); ok {
		r0 = rf(ctx, db, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]task_entity.Task)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, lib.QueryExecer, []int) error); ok {
		r1 = rf(ctx, db, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TaskRepo_GetTasksByIds_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTasksByIds'
type TaskRepo_GetTasksByIds_Call struct {
	*mock.Call
}

// GetTasksByIds is a helper method to define mock.On call
//   - ctx context.Context
//   - db lib.QueryExecer
//   - ids []int
func (_e *TaskRepo_Expecter) GetTasksByIds(ctx interface{}, db interface{}, ids interface{}) *TaskRepo_GetTasksByIds_Call {
	return &TaskRepo_GetTasksByIds_Call{Call: _e.mock.On("GetTasksByIds", ctx, db, ids)}
}

func (_c *TaskRepo_GetTasksByIds_Call) Run(run func(ctx context.Context, db lib.QueryExecer, ids []int)) *TaskRepo_GetTasksByIds_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(lib.QueryExecer), args[2].([]int))
	})
	return _c
}

func (_c *TaskRepo_GetTasksByIds_Call) Return(_a0 []task_entity.Task, _a1 error) *TaskRepo_GetTasksByIds_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TaskRepo_GetTasksByIds_Call) RunAndReturn(run func(context.Context, lib.QueryExecer, []int) ([]task_entity.Task, error)) *TaskRepo_GetTasksByIds_Call {
	_c.Call.Return(run)
	return _c
}

// GetTasksWithFilter provides a mock function with given fields: ctx, db, limit, offset
func (_m *TaskRepo) GetTasksWithFilter(ctx context.Context, db lib.QueryExecer, limit int, offset int) ([]task_entity.Task, error) {
	ret := _m.Called(ctx, db, limit, offset)

	var r0 []task_entity.Task
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, lib.QueryExecer, int, int) ([]task_entity.Task, error)); ok {
		return rf(ctx, db, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, lib.QueryExecer, int, int) []task_entity.Task); ok {
		r0 = rf(ctx, db, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]task_entity.Task)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, lib.QueryExecer, int, int) error); ok {
		r1 = rf(ctx, db, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TaskRepo_GetTasksWithFilter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTasksWithFilter'
type TaskRepo_GetTasksWithFilter_Call struct {
	*mock.Call
}

// GetTasksWithFilter is a helper method to define mock.On call
//   - ctx context.Context
//   - db lib.QueryExecer
//   - limit int
//   - offset int
func (_e *TaskRepo_Expecter) GetTasksWithFilter(ctx interface{}, db interface{}, limit interface{}, offset interface{}) *TaskRepo_GetTasksWithFilter_Call {
	return &TaskRepo_GetTasksWithFilter_Call{Call: _e.mock.On("GetTasksWithFilter", ctx, db, limit, offset)}
}

func (_c *TaskRepo_GetTasksWithFilter_Call) Run(run func(ctx context.Context, db lib.QueryExecer, limit int, offset int)) *TaskRepo_GetTasksWithFilter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(lib.QueryExecer), args[2].(int), args[3].(int))
	})
	return _c
}

func (_c *TaskRepo_GetTasksWithFilter_Call) Return(_a0 []task_entity.Task, _a1 error) *TaskRepo_GetTasksWithFilter_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TaskRepo_GetTasksWithFilter_Call) RunAndReturn(run func(context.Context, lib.QueryExecer, int, int) ([]task_entity.Task, error)) *TaskRepo_GetTasksWithFilter_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateTask provides a mock function with given fields: ctx, db, task
func (_m *TaskRepo) UpdateTask(ctx context.Context, db lib.QueryExecer, task task_entity.Task) error {
	ret := _m.Called(ctx, db, task)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, lib.QueryExecer, task_entity.Task) error); ok {
		r0 = rf(ctx, db, task)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TaskRepo_UpdateTask_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateTask'
type TaskRepo_UpdateTask_Call struct {
	*mock.Call
}

// UpdateTask is a helper method to define mock.On call
//   - ctx context.Context
//   - db lib.QueryExecer
//   - task task_entity.Task
func (_e *TaskRepo_Expecter) UpdateTask(ctx interface{}, db interface{}, task interface{}) *TaskRepo_UpdateTask_Call {
	return &TaskRepo_UpdateTask_Call{Call: _e.mock.On("UpdateTask", ctx, db, task)}
}

func (_c *TaskRepo_UpdateTask_Call) Run(run func(ctx context.Context, db lib.QueryExecer, task task_entity.Task)) *TaskRepo_UpdateTask_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(lib.QueryExecer), args[2].(task_entity.Task))
	})
	return _c
}

func (_c *TaskRepo_UpdateTask_Call) Return(_a0 error) *TaskRepo_UpdateTask_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TaskRepo_UpdateTask_Call) RunAndReturn(run func(context.Context, lib.QueryExecer, task_entity.Task) error) *TaskRepo_UpdateTask_Call {
	_c.Call.Return(run)
	return _c
}

// NewTaskRepo creates a new instance of TaskRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTaskRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *TaskRepo {
	mock := &TaskRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
