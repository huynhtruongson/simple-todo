package task_biz

import (
	"context"
	"testing"

	mock_db "github.com/sondev/todo-list/mock"
	task_entity "github.com/sondev/todo-list/services/task/entity"
	task_mock "github.com/sondev/todo-list/services/task/mock/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteTaskBiz_DeleteTask(t *testing.T) {
	ctx := context.Background()
	taskRepo := &task_mock.MockTaskRepo{}
	db := &mock_db.MockDB{}
	tx := &mock_db.MockTx{}
	tests := []struct {
		name      string
		taskId    int
		mock      func()
		expectErr error
	}{
		{
			name:   "should Delete task successfully",
			taskId: 1,
			mock: func() {
				taskRepo.On("GetTasksByIds", ctx, db, []int{1}).Once().Return([]task_entity.Task{{TaskID: 1}}, nil)
				db.On("BeginTx", ctx, mock.Anything).Once().Return(tx, nil)
				taskRepo.On("DeleteTask", ctx, tx, 1).Once().Return(nil)
				tx.On("Commit", ctx).Once().Return(nil)
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			biz := NewDeleteTaskBiz(db, taskRepo)
			tt.mock()
			err := biz.DeleteTask(ctx, tt.taskId)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
		})
	}
}
