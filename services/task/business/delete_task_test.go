package task_biz

import (
	"context"
	"testing"

	mock_db "github.com/huynhtruongson/simple-todo/mocks/lib"
	mock_repo "github.com/huynhtruongson/simple-todo/mocks/task"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteTaskBiz_DeleteTask(t *testing.T) {
	ctx := context.Background()
	taskRepo := mock_repo.NewTaskRepo(t)
	db := mock_db.NewDB(t)
	tx := mock_db.NewTx(t)
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
				taskRepo.EXPECT().GetTasksByIds(ctx, db, []int{1}).Once().Return([]task_entity.Task{{TaskID: 1}}, nil)
				db.EXPECT().BeginTx(ctx, mock.Anything).Once().Return(tx, nil)
				taskRepo.EXPECT().DeleteTask(ctx, tx, 1).Once().Return(nil)
				tx.EXPECT().Commit(ctx).Once().Return(nil)
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
