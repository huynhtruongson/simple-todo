package task_biz

import (
	"context"
	"testing"

	"github.com/sondev/todo-list/common"
	mock_db "github.com/sondev/todo-list/mock"
	task_entity "github.com/sondev/todo-list/services/task/entity"
	task_mock "github.com/sondev/todo-list/services/task/mock/repo"
	"github.com/stretchr/testify/assert"
)

func TestListTaskBiz_ListTask(t *testing.T) {
	ctx := context.Background()
	taskRepo := &task_mock.MockTaskRepo{}
	db := &mock_db.MockDB{}
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
			expectRes: TasksWithPaging{
				Tasks: []task_entity.Task{mockTask},
				Paging: common.Paging{
					Page:  2,
					Limit: 10,
					Total: 11,
				},
			},
			mock: func() {
				taskRepo.On("CountTask", ctx, db).Once().Return(11, nil)
				taskRepo.On("GetTasksWithFilter", ctx, db, 10, 10).Once().Return([]task_entity.Task{mockTask}, nil)
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
			expectRes: TasksWithPaging{
				Tasks: []task_entity.Task{mockTask},
				Paging: common.Paging{
					Page:  1,
					Limit: 10,
					Total: 1,
				},
			},
			mock: func() {
				taskRepo.On("CountTask", ctx, db).Once().Return(1, nil)
				taskRepo.On("GetTasksWithFilter", ctx, db, 10, 0).Once().Return([]task_entity.Task{mockTask}, nil)
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			biz := NewListTaskBiz(db, taskRepo)
			tt.mock()
			res, err := biz.ListTask(ctx, tt.paging)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
			assert.Equal(t, tt.expectRes.Tasks, res.Tasks)
			assert.Equal(t, tt.expectRes.Paging, res.Paging)
		})
	}
}
