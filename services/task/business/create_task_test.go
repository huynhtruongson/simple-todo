package task_biz

import (
	"context"
	"testing"

	"github.com/huynhtruongson/simple-todo/common"
	mock_db "github.com/huynhtruongson/simple-todo/mocks/lib"
	mock_repo "github.com/huynhtruongson/simple-todo/mocks/task"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTaskBiz_CreateTask(t *testing.T) {
	ctx := context.Background()
	userRepo := mock_repo.NewUserRepo(t)
	taskRepo := mock_repo.NewTaskRepo(t)
	db := mock_db.NewDB(t)
	tx := mock_db.NewTx(t)
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
				userRepo.EXPECT().GetUsersByUserIds(ctx, db, []int{1}).Once().Return([]user_entity.User{{UserID: 1}}, nil)
				db.EXPECT().BeginTx(ctx, mock.Anything).Once().Return(tx, nil)
				taskRepo.EXPECT().CreateTask(ctx, tx, task_entity.Task{
					Title:  "title",
					UserID: 1,
					Status: 1,
				}).Once().Return(1, nil)
				tx.EXPECT().Commit(ctx).Once().Return(nil)
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
	userRepo := mock_repo.NewUserRepo(t)
	db := mock_db.NewDB(t)
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
				userRepo.EXPECT().GetUsersByUserIds(ctx, db, []int{1}).Once().Return([]user_entity.User{}, nil)
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
