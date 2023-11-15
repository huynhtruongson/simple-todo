package task_biz

import (
	"context"
	"testing"

	"github.com/huynhtruongson/simple-todo/common"
	mock_db "github.com/huynhtruongson/simple-todo/mocks/lib"
	mock_repo "github.com/huynhtruongson/simple-todo/mocks/task"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"

	"github.com/stretchr/testify/assert"
)

func TestListTaskBiz_ListTask(t *testing.T) {
	ctx := context.Background()
	taskRepo := mock_repo.NewTaskRepo(t)
	db := mock_db.NewDB(t)
	mockTask := task_entity.Task{
		TaskID: 1,
		UserID: 1,
		Title:  "title",
	}
	tests := []struct {
		name      string
		taskId    int
		mock      func()
		paging    common.Paging
		filter    common.Filter
		expectRes TasksWithPaging
		expectErr error
	}{
		{
			name:   "should List task successfully",
			taskId: 1,
			paging: common.Paging{
				Page:  2,
				Limit: 10,
			},
			filter: common.Filter{
				UserID: 1,
			},
			expectRes: TasksWithPaging{
				Tasks: []task_entity.Task{mockTask},
				Paging: common.Paging{
					Page:  2,
					Limit: 10,
					Total: 11,
				},
			},
			mock: func() {
				taskRepo.EXPECT().CountTask(ctx, db, 1).Once().Return(11, nil)
				taskRepo.EXPECT().GetTasksWithFilter(ctx, db, 1, 10, 10).Once().Return([]task_entity.Task{mockTask}, nil)
			},
			expectErr: nil,
		},
		{
			name:   "should List task successfully when paging params is invalid",
			taskId: 1,
			paging: common.Paging{
				Page:  0,
				Limit: 0,
			},
			filter: common.Filter{
				UserID: 1,
			},
			expectRes: TasksWithPaging{
				Tasks: []task_entity.Task{mockTask},
				Paging: common.Paging{
					Page:  1,
					Limit: 10,
					Total: 1,
				},
			},
			mock: func() {
				taskRepo.EXPECT().CountTask(ctx, db, 1).Once().Return(1, nil)
				taskRepo.EXPECT().GetTasksWithFilter(ctx, db, 1, 10, 0).Once().Return([]task_entity.Task{mockTask}, nil)
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			biz := NewListTaskBiz(db, taskRepo)
			tt.mock()
			res, err := biz.ListTask(ctx, tt.paging, tt.filter)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
			assert.Equal(t, tt.expectRes.Tasks, res.Tasks)
			assert.Equal(t, tt.expectRes.Paging, res.Paging)
		})
	}
}
