package task_service

import (
	"context"
	"testing"

	"github.com/huynhtruongson/simple-todo/common"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"

	"github.com/stretchr/testify/assert"
)

func TestListTaskBiz_ListTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockTask := task_entity.Task{
		TaskID: 1,
		UserID: 1,
		Title:  "title",
	}
	tests := []struct {
		name      string
		taskId    int
		mock      func(prop *MockServiceProp)
		paging    common.Paging
		filter    common.Filter
		expectRes task_entity.TasksWithPaging
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
			expectRes: task_entity.TasksWithPaging{
				Tasks: []task_entity.Task{mockTask},
				Paging: common.Paging{
					Page:  2,
					Limit: 10,
					Total: 11,
				},
			},
			mock: func(prop *MockServiceProp) {
				prop.TaskRepo.EXPECT().CountTask(ctx, prop.DB, 1).Once().Return(11, nil)
				prop.TaskRepo.EXPECT().GetTasksWithFilter(ctx, prop.DB, 1, 10, 10).Once().Return([]task_entity.Task{mockTask}, nil)
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
			expectRes: task_entity.TasksWithPaging{
				Tasks: []task_entity.Task{mockTask},
				Paging: common.Paging{
					Page:  1,
					Limit: 10,
					Total: 1,
				},
			},
			mock: func(prop *MockServiceProp) {
				prop.TaskRepo.EXPECT().CountTask(ctx, prop.DB, 1).Once().Return(1, nil)
				prop.TaskRepo.EXPECT().GetTasksWithFilter(ctx, prop.DB, 1, 10, 0).Once().Return([]task_entity.Task{mockTask}, nil)
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sv, prop := NewMockTaskService(t)
			tt.mock(prop)
			res, err := sv.ListTask(ctx, tt.paging, tt.filter)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
			assert.Equal(t, tt.expectRes.Tasks, res.Tasks)
			assert.Equal(t, tt.expectRes.Paging, res.Paging)
		})
	}
}
