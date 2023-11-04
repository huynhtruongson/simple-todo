package task_biz

import (
	"context"
	"testing"

	mock_db "github.com/sondev/todo-list/mock"
	task_entity "github.com/sondev/todo-list/services/task/entity"
	task_mock "github.com/sondev/todo-list/services/task/mock/repo"
	user_entity "github.com/sondev/todo-list/services/user/entity"
	user_mock "github.com/sondev/todo-list/services/user/mock/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateTaskBiz_UpdateTask(t *testing.T) {
	ctx := context.Background()
	userRepo := &user_mock.MockUserRepo{}
	taskRepo := &task_mock.MockTaskRepo{}
	db := &mock_db.MockDB{}
	tx := &mock_db.MockTx{}
	tests := []struct {
		name      string
		task      task_entity.Task
		mock      func()
		expectErr error
	}{
		{
			name: "should update task successfully",
			task: task_entity.Task{
				TaskID: 1,
				Title:  "title",
				UserID: 1,
				Status: 1,
			},
			mock: func() {
				taskRepo.On("GetTasksByIds", ctx, db, []int{1}).Once().Return([]task_entity.Task{{TaskID: 1}}, nil)
				userRepo.On("GetUsersByUserIds", ctx, db, []int{1}).Once().Return([]user_entity.User{{UserID: 1}}, nil)
				db.On("BeginTx", ctx, mock.Anything).Once().Return(tx, nil)
				taskRepo.On("UpdateTask", ctx, tx, task_entity.Task{
					TaskID: 1,
					Title:  "title",
					UserID: 1,
					Status: 1,
				}).Once().Return(nil)
				tx.On("Commit", ctx).Once().Return(nil)
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			biz := NewUpdateTaskBiz(db, taskRepo, userRepo)
			tt.mock()
			err := biz.UpdateTask(ctx, tt.task)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
		})
	}
}
