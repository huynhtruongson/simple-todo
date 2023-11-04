package task_biz

import (
	"context"
	"testing"

	"github.com/sondev/todo-list/common"
	mock_db "github.com/sondev/todo-list/mock"
	task_entity "github.com/sondev/todo-list/services/task/entity"
	task_mock "github.com/sondev/todo-list/services/task/mock/repo"
	user_entity "github.com/sondev/todo-list/services/user/entity"
	user_mock "github.com/sondev/todo-list/services/user/mock/repo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTaskBiz_CreateTask(t *testing.T) {
	ctx := context.Background()
	userRepo := &user_mock.MockUserRepo{}
	taskRepo := &task_mock.MockTaskRepo{}
	db := &mock_db.MockDB{}
	tx := &mock_db.MockTx{}
	tests := []struct {
		name         string
		task         task_entity.Task
		mock         func()
		expectErr    error
		expectTaskID int
	}{
		{
			name: "should return taskID when create task successfully",
			task: task_entity.Task{
				Title:  "title",
				UserID: 1,
				Status: 1,
			},
			mock: func() {
				userRepo.On("GetUsersByUserIds", ctx, db, []int{1}).Once().Return([]user_entity.User{{UserID: 1}}, nil)
				db.On("BeginTx", ctx, mock.Anything).Once().Return(tx, nil)
				taskRepo.On("CreateTask", ctx, tx, task_entity.Task{
					Title:  "title",
					UserID: 1,
					Status: 1,
				}).Once().Return(1, nil)
				tx.On("Commit", ctx).Once().Return(nil)
			},
			expectErr:    nil,
			expectTaskID: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			biz := NewCreateTaskBiz(db, taskRepo, userRepo)
			tt.mock()
			taskID, err := biz.CreateTask(ctx, tt.task)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
			assert.Equal(t, tt.expectTaskID, taskID)
		})
	}
}

func TestCreateTaskBiz_ValidateTask(t *testing.T) {
	ctx := context.Background()
	userRepo := &user_mock.MockUserRepo{}
	db := &mock_db.MockDB{}
	tests := []struct {
		name      string
		task      task_entity.Task
		mock      func()
		expectErr *common.AppError
	}{
		{
			name: "should throw error when title is empty",
			task: task_entity.Task{
				Title: "",
			},
			mock:      func() {},
			expectErr: common.NewInvalidRequestError(nil, task_entity.ErrorTitleIsEmpty, "ValidateTask"),
		},
		{
			name: "should throw error when userID is empty",
			task: task_entity.Task{
				Title: "title",
			},
			mock:      func() {},
			expectErr: common.NewInvalidRequestError(nil, task_entity.ErrorUserIsEmpty, "ValidateTask"),
		},
		{
			name: "should throw error when status is invalid",
			task: task_entity.Task{
				Title:  "title",
				UserID: 1,
				Status: 3,
			},
			mock:      func() {},
			expectErr: common.NewInvalidRequestError(nil, task_entity.ErrorInvalidStatus, "ValidateTask"),
		},
		{
			name: "should throw error when userID does not exist",
			task: task_entity.Task{
				Title:  "title",
				UserID: 1,
				Status: 1,
			},
			mock: func() {
				userRepo.On("GetUsersByUserIds", ctx, db, []int{1}).Once().Return([]user_entity.User{}, nil)
			},
			expectErr: common.NewInvalidRequestError(nil, task_entity.ErrorUserNotFound, "ValidateTask"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			biz := CreateTaskBiz{
				DB:       db,
				UserRepo: userRepo,
			}
			tt.mock()
			err := biz.ValidateTask(ctx, tt.task)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr.Code, err.(*common.AppError).Code)
				assert.Equal(t, tt.expectErr.Message, err.(*common.AppError).Message)
			}
		})
	}
}
