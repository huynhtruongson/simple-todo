package task_service

import (
	"context"
	"testing"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/field"
	mock_db "github.com/huynhtruongson/simple-todo/mocks/lib"
	mock_repo "github.com/huynhtruongson/simple-todo/mocks/task"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockServiceProp struct {
	DB       *mock_db.DB
	TX       *mock_db.Tx
	TaskRepo *mock_repo.TaskRepo
	UserRepo *mock_repo.UserRepo
}

func NewMockTaskService(t *testing.T) (*TaskService, *MockServiceProp) {
	userRepo := mock_repo.NewUserRepo(t)
	taskRepo := mock_repo.NewTaskRepo(t)
	db := mock_db.NewDB(t)
	tx := mock_db.NewTx(t)
	return &TaskService{
			DB:       db,
			TaskRepo: taskRepo,
			UserRepo: userRepo,
		}, &MockServiceProp{
			DB:       db,
			TX:       tx,
			TaskRepo: taskRepo,
			UserRepo: userRepo,
		}
}
func TestCreateTaskBiz_CreateTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tests := []struct {
		name         string
		task         task_entity.Task
		mock         func(prop *MockServiceProp)
		expectErr    error
		expectTaskID int
	}{
		{
			name: "should return taskID when create task successfully",
			task: task_entity.Task{
				Title:  field.NewString("title"),
				UserID: field.NewInt(1),
				Status: 1,
			},
			mock: func(prop *MockServiceProp) {
				prop.UserRepo.EXPECT().GetUsersByUserIds(ctx, prop.DB, []int{1}).Once().Return([]user_entity.User{{UserID: field.NewInt(1)}}, nil)
				prop.DB.EXPECT().BeginTx(ctx, mock.Anything).Once().Return(prop.TX, nil)
				prop.TaskRepo.EXPECT().CreateTask(ctx, prop.TX, task_entity.Task{
					Title:  field.NewString("title"),
					UserID: field.NewInt(1),
					Status: 1,
				}).Once().Return(1, nil)
				prop.TX.EXPECT().Commit(ctx).Once().Return(nil)
			},
			expectErr:    nil,
			expectTaskID: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sv, prop := NewMockTaskService(t)
			tt.mock(prop)
			taskID, err := sv.CreateTask(ctx, tt.task)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
			assert.Equal(t, tt.expectTaskID, taskID)
		})
	}
}

func TestCreateTaskBiz_ValidateTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tests := []struct {
		name      string
		task      task_entity.Task
		mock      func(prop *MockServiceProp)
		expectErr *common.AppError
	}{
		{
			name: "should throw error when title is empty",
			task: task_entity.Task{
				Title: field.NewNullString(),
			},
			mock:      func(prop *MockServiceProp) {},
			expectErr: common.NewInvalidRequestError(task_entity.ErrorTitleIsEmpty, task_entity.ErrorTitleIsEmpty.Error(), "ValidateTask"),
		},
		{
			name: "should throw error when userID is empty",
			task: task_entity.Task{
				Title: field.NewString("title"),
			},
			mock:      func(prop *MockServiceProp) {},
			expectErr: common.NewInvalidRequestError(task_entity.ErrorUserIsEmpty, task_entity.ErrorUserIsEmpty.Error(), "ValidateTask"),
		},
		{
			name: "should throw error when status is invalid",
			task: task_entity.Task{
				Title:  field.NewString("title"),
				UserID: field.NewInt(1),
				Status: 3,
			},
			mock:      func(prop *MockServiceProp) {},
			expectErr: common.NewInvalidRequestError(task_entity.ErrorInvalidStatus, task_entity.ErrorInvalidStatus.Error(), "ValidateTask"),
		},
		{
			name: "should throw error when userID does not exist",
			task: task_entity.Task{
				Title:  field.NewString("title"),
				UserID: field.NewInt(1),
				Status: 1,
			},
			mock: func(prop *MockServiceProp) {
				prop.UserRepo.EXPECT().GetUsersByUserIds(ctx, prop.DB, []int{1}).Once().Return([]user_entity.User{}, nil)
			},
			expectErr: common.NewInvalidRequestError(task_entity.ErrorUserNotFound, task_entity.ErrorUserNotFound.Error(), "ValidateTask"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sv, prop := NewMockTaskService(t)
			tt.mock(prop)
			err := sv.ValidateTask(ctx, tt.task)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr.Code, err.(*common.AppError).Code)
				assert.Equal(t, tt.expectErr.Message, err.(*common.AppError).Message)
			}
		})
	}
}
