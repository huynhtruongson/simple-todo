package task_biz

import (
	"context"
	"testing"

	mock_db "github.com/huynhtruongson/simple-todo/mocks/lib"
	mock_repo "github.com/huynhtruongson/simple-todo/mocks/task"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateTaskBiz_UpdateTask(t *testing.T) {
	ctx := context.Background()
	userRepo := mock_repo.NewUserRepo(t)
	taskRepo := mock_repo.NewTaskRepo(t)
	db := mock_db.NewDB(t)
	tx := mock_db.NewTx(t)
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
				taskRepo.EXPECT().GetTasksByIds(ctx, db, 1, []int{1}).Once().Return([]task_entity.Task{{TaskID: 1}}, nil)
				userRepo.EXPECT().GetUsersByUserIds(ctx, db, []int{1}).Once().Return([]user_entity.User{{UserID: 1}}, nil)
				db.EXPECT().BeginTx(ctx, mock.Anything).Once().Return(tx, nil)
				taskRepo.EXPECT().UpdateTask(ctx, tx, task_entity.Task{
					TaskID: 1,
					Title:  "title",
					UserID: 1,
					Status: 1,
				}).Once().Return(nil)
				tx.EXPECT().Commit(ctx).Once().Return(nil)
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
