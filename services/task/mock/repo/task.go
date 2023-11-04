package task_mock

import (
	"context"

	"github.com/sondev/todo-list/lib"
	task_entity "github.com/sondev/todo-list/services/task/entity"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepo struct {
	mock.Mock
}

func (m *MockTaskRepo) CreateTask(ctx context.Context, db lib.QueryExecer, task task_entity.Task) (int, error) {
	args := m.Called(ctx, db, task)
	return args.Int(0), args.Error(1)
}
func (m *MockTaskRepo) UpdateTask(ctx context.Context, db lib.QueryExecer, task task_entity.Task) error {
	args := m.Called(ctx, db, task)
	return args.Error(0)
}
func (m *MockTaskRepo) GetTasksByIds(ctx context.Context, db lib.QueryExecer, ids []int) ([]task_entity.Task, error) {
	args := m.Called(ctx, db, ids)
	return args.Get(0).([]task_entity.Task), args.Error(1)
}

func (m *MockTaskRepo) DeleteTask(ctx context.Context, db lib.QueryExecer, id int) error {
	args := m.Called(ctx, db, id)
	return args.Error(0)
}

func (m *MockTaskRepo) CountTask(ctx context.Context, db lib.QueryExecer) (int, error) {
	args := m.Called(ctx, db)
	return args.Int(0), args.Error(1)
}

func (m *MockTaskRepo) GetTasksWithFilter(ctx context.Context, db lib.QueryExecer, limit, offset int) ([]task_entity.Task, error) {
	args := m.Called(ctx, db, limit, offset)
	return args.Get(0).([]task_entity.Task), args.Error(1)
}
