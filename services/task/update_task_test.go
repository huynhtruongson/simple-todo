package task_service

import (
	"context"
	"testing"

	"github.com/huynhtruongson/simple-todo/field"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateTaskBiz_UpdateTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tests := []struct {
		name      string
		task      task_entity.Task
		mock      func(prop *MockServiceProp)
		expectErr error
	}{
		{
			name: "should update task successfully",
			task: task_entity.Task{
				TaskID: field.NewInt(1),
				Title:  field.NewString("title"),
				UserID: field.NewInt(1),
				Status: 1,
			},
			mock: func(prop *MockServiceProp) {
				prop.TaskRepo.EXPECT().GetTasksByIds(ctx, prop.DB, 1, []int{1}).Once().Return([]task_entity.Task{{TaskID: field.NewInt(1)}}, nil)
				prop.UserRepo.EXPECT().GetUsersByUserIds(ctx, prop.DB, []int{1}).Once().Return([]user_entity.User{{UserID: field.NewInt(1)}}, nil)
				prop.DB.EXPECT().BeginTx(ctx, mock.Anything).Once().Return(prop.TX, nil)
				prop.TaskRepo.EXPECT().UpdateTask(ctx, prop.TX, task_entity.Task{
					TaskID: field.NewInt(1),
					Title:  field.NewString("title"),
					UserID: field.NewInt(1),
					Status: 1,
				}).Once().Return(nil)
				prop.TX.EXPECT().Commit(ctx).Once().Return(nil)
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sv, prop := NewMockTaskService(t)
			tt.mock(prop)
			err := sv.UpdateTask(ctx, tt.task)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr, err)
				return
			}
		})
	}
}
